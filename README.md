# unchain ⛓️

Release tool with automatic changelog generation and next SemVer version calculation based on conventional commits.

## Install

```sh
go get github.com/flume/release-version
```

## Usage

See help for available commands and parameters.

```sh
unchain -h
```

## Commands

### release

Run in your terminal:

```sh
unchain release
```

- `dir` optional *(default: workdir)*, directory
- `suppressPush` optional *(default: false)*, suppress the git push

### Starting from some previously defined version
Add a commit with a first line that includes the previous version that looks like this:
`chore(release): 1.7.1`

The next run will then account for any commits between that chore commit and the current HEAD

#### How It Works

Automatically detects the last version and bumps the `patch`, `minor` or `major` semver component based on the conventional commits since that release.
If there is no commit found related to previous version it will release `1.0.0`.

*What It Does*

* Detects the next SemVer version based on commit history
* Detects the previous version from release commits made by this tool or from package.json if any
* Creates or prepends `CHANGELOG.md`
* *(optional)* Execs `npm version` if finds package.json
* Git tags release
* *(optional)* `npm publish` if finds package.json
* Runs `git push` to sync with remote

*CHANGELOG.md example*

```
### [1.7.1](https://github.com/flume/thing/compare/1.7.0...1.7.1) (2020-01-22)

### Bug Fixes

* **k8s:** remove lonely config params ([b6d2547](https://github.com/flume/proxy/commit/b6d254762cec7bcf42cbedc0d0ea41d24331dca0))

```

*Commits example*
- (always): chore(release): 1.0.0

*Tag created*

- `1.0.0` (with package.json, v1.0.0)

Skips non API facing commits from the changelog like `test`, `chore` and `refactor`.

## semver

Detect SemVer change since latest Git Tag.

```sh
$ unchain semver
Change Detected: major
Next Version: 1.0.0
```

## parse

Parse commits since latest Git Tag.

```sh
$ unchain parse
hash,semver,type,component,description,body,footer
ecd94da5b9f10c04ce53723729ae7068cc73557e,major,fix,foo,fifth commit,body,BREAKING CHANGE: so breaking
29afc9699602e73418395226f22389a5271c5e58,major,fix,bar,fourth commit,BREAKING CHANGE: blabla,
6289d27b800d3966ec7f14394ff4c48b08dd5976,patch,fix,foo,third commit,body,
998df6abedeeb0e090986b5de3a89e62c03c436d,patch,chore,foo,second commit,,
a4a95856d51dc3018170f2a854581590d1a27687,minor,feat,foo,initial commit,,
```

## Background

Fork of https://github.com/hekike/unchain

Follows:

- https://semver.org
- https://www.conventionalcommits.org
