package git

import (
	"errors"
	"fmt"
	"strings"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Release Git release
func Release(dir string, version string, user *User, suppressPush bool) error {
	_, err := Tag(dir, version, user)
	if err != nil {
		return fmt.Errorf("[Release] tag: %v", err)
	}

	if !suppressPush {
		err = Push(dir, version)
		if err != nil {
			return fmt.Errorf("[Release] push: %v", err)
		}
	}

	return nil
}

// Commit commit file
func Commit(dir string, file string, message string, user *User, branchName string) (err error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("[Git] open repo: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("[Git] worktree: %v", err)
	}

	headRef, err := r.Head()
	if err != nil {
		return fmt.Errorf("[Git] head ref: %v", err)
	}

	defer func() {
		cerr := w.Checkout(&git.CheckoutOptions{
			Branch: headRef.Name(),
		})
		if cerr != nil {
			err = fmt.Errorf("%v: %v", cerr, err)
		}
	}()

	if branchName != "" {
		// Check out the desired branch
		err = r.Fetch(&git.FetchOptions{})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return fmt.Errorf("[Git] fetch: %v", err)
		}

		// Try to get the reference of the desired branch
		var branchRef *plumbing.Reference
		branchRef, err := r.Reference(plumbing.NewRemoteReferenceName("origin", branchName), false)
		if err != nil {
			return fmt.Errorf("[Git] get remote branch reference: %v", err)
		}

		if branchRef == nil {
			return fmt.Errorf("[Git] branch '%s' not found", branchName)
		}

		// Checkout the branch by creating a new local branch based on the remote branch commit
		err = w.Checkout(&git.CheckoutOptions{
			Hash:   branchRef.Hash(),
			Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName)),
			Create: true,
			Keep:   true,
		})
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				err = w.Checkout(&git.CheckoutOptions{
					Branch: plumbing.NewBranchReferenceName(branchName),
					Keep:   true,
				})
				if err != nil {
					return fmt.Errorf("[Git] checkout local branch: %v", err)
				}
			}
		}
	}

	_, err = w.Add(file)
	if err != nil {
		return fmt.Errorf("[Git] worktree add (%s): %v", file, err)
	}

	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  user.Name,
			Email: user.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("[Git] commit: %v", err)
	}

	return nil
}

// Tag tag last commit
func Tag(dir string, version string, user *User) (*plumbing.Reference, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("[Tag] open repo: %v", err)
	}

	head, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("[Tag] head: %v", err)
	}

	ref, err := r.CreateTag(version, head.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  user.Name,
			Email: user.Email,
			When:  time.Now(),
		},
		Message: version,
	})
	if err != nil {
		return nil, fmt.Errorf("[Tag] create git tag: %v", err)
	}

	return ref, nil
}

func GetLatestTag(dir string) (string, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return "", fmt.Errorf("[Tag] open repo: %v", err)
	}

	tags, err := r.Tags()
	if err != nil {
		return "", err
	}

	var latestTag string
	var latestCommit *object.Commit

	err = tags.ForEach(func(t *plumbing.Reference) error {
		// We're interested in annotated tags only
		if t.Name().IsTag() {
			tagObj, err := r.TagObject(t.Hash())
			if err != nil {
				switch {
				case errors.Is(err, plumbing.ErrObjectNotFound):
					return nil
				default:
					return fmt.Errorf("tag object: %v", err)
				}
			}

			// Get the commit object associated with the tag
			commitObj, err := r.CommitObject(tagObj.Target)
			if err != nil {
				return fmt.Errorf("commit object: %v", err)
			}

			// Check if this tag is the latest one
			if latestCommit == nil || commitObj.Committer.When.After(latestCommit.Committer.When) {
				latestCommit = commitObj
				latestTag = t.Name().Short()
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return latestTag, nil
}

// Push push to remote
func Push(dir string, version string) error {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("[Push] open repo: %v", err)
	}

	// push using default options
	tagRef := fmt.Sprintf("refs/tags/%s:refs/tags/%s", version, version)
	err = r.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/heads/master:refs/heads/master"),
			config.RefSpec(tagRef),
		},
	})
	if err != nil {
		switch err {
		case git.ErrRemoteNotFound:
			return nil
		default:
			return fmt.Errorf("[Push] push: %v", err)
		}
	}

	return nil
}

func GetRemotePath(dir string) (string, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return "", fmt.Errorf("[Remote Path] open repo: %v", err)
	}

	remotes, err := r.Remotes()
	if err != nil {
		return "", fmt.Errorf("[Remote Path] get remotes: %v", err)
	}

	var path string
	for _, r := range remotes {
		if !strings.Contains(r.String(), "origin") {
			continue
		}
		urls := r.Config().URLs
		if len(urls) < 1 {
			return "", fmt.Errorf("[Remote Path] couldn't find remote urls")
		}
		path = urls[0]
		break
	}

	if path != "" {
		splt := strings.Split(path, "@")
		// take the second half and use it in the path to return
		if len(splt) > 1 {
			path = splt[len(splt)-1]
		}
	}

	path = strings.TrimSuffix(path, ".git")
	path = strings.Replace(path, ":", "/", -1)
	return path, nil
}
