package changelog

/*
func TestGenerate(t *testing.T) {
	commits := []parser.ConventionalCommit{}

	date := time.Now().Format("2006-01-02")

	// Zero commits
	res := Generate("1.0.0", commits)
	assert.Equal(t, "<a name=\"1.0.0\"></a>\n## 1.0.0 ("+date+")\n"+
		"\n\n* There is no user facing commit in this version"+
		"\n"+
		"\n\n", res)

	// Patch only
	commits = append(commits, parser.ConventionalCommit{
		Type:         "fix",
		Component:    "foo",
		Description:  "fixing dividing by zero",
		SemVerChange: parser.Patch,
	})
	res = Generate("1.0.0", commits)
	assert.Equal(t, "<a name=\"1.0.0\"></a>\n## 1.0.0 ("+date+")\n"+
		"\n\n#### Bug Fixes\n"+
		"\n* **foo:** fixing dividing by zero "+
		"\n"+
		"\n\n", res)

	// Patch and minor
	commits = append(commits, parser.ConventionalCommit{
		Type:         "fix",
		Component:    "foo",
		Description:  "add new option",
		SemVerChange: parser.Minor,
	})
	res = Generate("1.0.0", commits)
	assert.Equal(t, "<a name=\"1.0.0\"></a>\n## 1.0.0 ("+date+")\n"+
		"\n\n#### Bug Fixes\n"+
		"\n* **foo:** fixing dividing by zero "+
		"\n\n#### Features\n"+
		"\n* **foo:** add new option "+
		"\n"+
		"\n\n", res)

	// Patch, minor and major
	commits = append(commits, parser.ConventionalCommit{
		Type:         "fix",
		Component:    "foo",
		Breaking:     "renaming input",
		SemVerChange: parser.Major,
	})
	res = Generate("1.0.0", commits)
	assert.Equal(t, "<a name=\"1.0.0\"></a>\n## 1.0.0 ("+date+")\n"+
		"\n\n#### Bug Fixes\n"+
		"\n* **foo:** fixing dividing by zero "+
		"\n\n#### Features\n"+
		"\n* **foo:** add new option "+
		"\n\n#### Breaking Changes\n"+
		"\n* renaming input "+
		"\n"+
		"\n\n", res)
}
*/
