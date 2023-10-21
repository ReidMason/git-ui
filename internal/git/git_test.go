package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type MockGitCommandRunner struct {
	diffResponse   string
	statusResponse string
}

func (g MockGitCommandRunner) Stage(filepath string)          {}
func (g MockGitCommandRunner) Unstage(filepath string)        {}
func (g MockGitCommandRunner) Commit(commitMessage string)    {}
func (g MockGitCommandRunner) GetDiff(filepath string) string { return g.diffResponse }
func (g MockGitCommandRunner) GetStatus() string              { return g.statusResponse }

func TestGetDiff(t *testing.T) {
	testCases := []struct {
		name, rawDiff  string
		expectedResult Diff
	}{
		{
			name: "Real world use case",
			rawDiff: `diff --git a/git-ui/testfile.txt b/git-ui/testfile.txt
index 35b5809..4492ac6 100644
--- a/git-ui/testfile.txt
+++ b/git-ui/testfile.txt
@@ -1,4 +1,5 @@
 This is a test file
-
 These lines are committed now
-I have added some more content
+I have added this is a change more content
+
+This is a new thing

 type model struct {
-       ldiff     string
-       rdiff     string
+       ldiff     []git.DiffLine
+       rdiff     []git.DiffLine
        lviewport viewport.Model
        rviewport viewport.Model
        ready     bool
 }`,
			expectedResult: Diff{
				Diff1: []DiffLine{
					{Content: "This is a test file", Type: Neutral},
					{Content: "", Type: Removal},
					{Content: "These lines are committed now", Type: Neutral},
					{Content: "I have added some more content", Type: Removal},
					{Content: "", Type: Blank},
					{Content: "", Type: Blank},
					{Content: "", Type: Neutral},
					{Content: "type model struct {", Type: Neutral},
					{Content: "       ldiff     string", Type: Removal},
					{Content: "       rdiff     string", Type: Removal},
					{Content: "       lviewport viewport.Model", Type: Neutral},
					{Content: "       rviewport viewport.Model", Type: Neutral},
					{Content: "       ready     bool", Type: Neutral},
					{Content: "}", Type: Neutral},
				},
				Diff2: []DiffLine{
					{Content: "This is a test file", Type: Neutral},
					{Content: "", Type: Blank},
					{Content: "These lines are committed now", Type: Neutral},
					{Content: "I have added this is a change more content", Type: Addition},
					{Content: "", Type: Addition},
					{Content: "This is a new thing", Type: Addition},
					{Content: "", Type: Neutral},
					{Content: "type model struct {", Type: Neutral},
					{Content: "       ldiff     []git.DiffLine", Type: Removal},
					{Content: "       rdiff     []git.DiffLine", Type: Removal},
					{Content: "       lviewport viewport.Model", Type: Neutral},
					{Content: "       rviewport viewport.Model", Type: Neutral},
					{Content: "       ready     bool", Type: Neutral},
					{Content: "}", Type: Neutral},
				},
			},
		},
		{
			name: "Trailing removal",
			rawDiff: `diff --git a/git-ui/testfile.txt b/git-ui/testfile.txt
index 35b5809..4492ac6 100644
--- a/git-ui/testfile.txt
+++ b/git-ui/testfile.txt
@@ -1,4 +1,5 @@
-This is a test file`,
			expectedResult: Diff{
				Diff1: []DiffLine{
					{Content: "This is a test file", Type: Removal},
				},
				Diff2: []DiffLine{
					{Content: "", Type: Blank},
				},
			},
		},
		{
			name: "Trailing addition",
			rawDiff: `diff --git a/git-ui/testfile.txt b/git-ui/testfile.txt
index 35b5809..4492ac6 100644
--- a/git-ui/testfile.txt
+++ b/git-ui/testfile.txt
@@ -1,4 +1,5 @@
+This is a test file`,
			expectedResult: Diff{
				Diff1: []DiffLine{
					{Content: "", Type: Blank},
				},
				Diff2: []DiffLine{
					{Content: "This is a test file", Type: Addition},
				},
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc

		commandRunner := MockGitCommandRunner{diffResponse: tc.rawDiff}
		git := New(commandRunner)
		result := git.GetDiff("path/to/file")

		for i, expectedDiffLine := range tc.expectedResult.Diff1 {
			resultDiffLine := result.Diff1[i]

			if !cmp.Equal(expectedDiffLine, resultDiffLine) {
				t.Fatalf("Diff1 failed test '%s' line %d.\nExpected: '%v'\n     Got: '%v'", tc.name, i+1, expectedDiffLine, resultDiffLine)
			}
		}

		for i, expectedDiffLine := range tc.expectedResult.Diff2 {
			resultDiffLine := result.Diff2[i]

			if expectedDiffLine.Content != resultDiffLine.Content {
				t.Fatalf("Diff2 failed test '%s' line %d.\nExpected: '%s'\n     Got: '%s'", tc.name, i+1, expectedDiffLine.Content, resultDiffLine.Content)
			}
		}
	}
}

func TestGetStatus(t *testing.T) {
	testCases := []struct {
		name      string
		rawStatus string
		expected  GitStatus
	}{
		{
			name: "Just status",
			rawStatus: `# branch.oid c86e7ed35f16570194c2308a2f8cb53155d0440d

# branch.head main
# branch.upstream origin/main
# branch.ab +0 -0
1 M. N... 100644 100644 100644 1cdd739f6591c3aca07eab977748142a1ba14056 c345bc6f17650da4f51350e8faa56e4f4c61663e main.go`,
			expected: func() GitStatus {
				rootDir := Directory{
					Name:   "Root",
					Parent: nil,
					Files:  []File{},
				}

				file := File{
					Name:           "main.go",
					Dirpath:        ".",
					Parent:         &rootDir,
					IndexStatus:    'M',
					WorktreeStatus: '.',
				}
				rootDir.Files = append(rootDir.Files, file)

				return GitStatus{
					Head:      "main",
					Upstream:  "origin/main",
					Directory: &rootDir,
				}
			}(),
		},
		{
			name: "Just status",
			rawStatus: `# branch.oid c86e7ed35f16570194c2308a2f8cb53155d0440d

# branch.head main
# branch.upstream origin/main
# branch.ab +0 -0
1 .. N... 100644 100644 100644 1cdd739f6591c3aca07eab977748142a1ba14056 c345bc6f17650da4f51350e8faa56e4f4c61663e Directory/Internal/main.go
1 MM N... 100644 100644 100644 1cdd739f6591c3aca07eab977748142a1ba14056 c345bc6f17650da4f51350e8faa56e4f4c61663e Directory/Internal/lib.go`,
			expected: func() GitStatus {
				rootDir := Directory{
					Name:   "Root",
					Parent: nil,
					Files:  []File{},
				}

				directory := Directory{
					Name:   "Directory",
					Parent: &rootDir,
					Files:  []File{},
				}
				rootDir.Directories = append(rootDir.Directories, &directory)

				internalDir := Directory{
					Name:   "Internal",
					Parent: &directory,
					Files:  []File{},
				}
				directory.Directories = append(directory.Directories, &internalDir)

				file := File{
					Name:           "main.go",
					Dirpath:        "Directory/Internal",
					Parent:         &internalDir,
					IndexStatus:    '.',
					WorktreeStatus: '.',
				}
				internalDir.Files = append(internalDir.Files, file)

				lib := File{
					Name:           "lib.go",
					Dirpath:        "Directory/Internal",
					Parent:         &internalDir,
					IndexStatus:    'M',
					WorktreeStatus: 'M',
				}
				internalDir.Files = append(internalDir.Files, lib)

				return GitStatus{
					Head:      "main",
					Upstream:  "origin/main",
					Directory: &rootDir,
				}
			}(),
		},
		{
			name: "Status without upstream",
			rawStatus: `# branch.oid 74700949a5ce67be9cb5ee97434df52846caec01
		# branch.head testing`,
			expected: func() GitStatus {
				return GitStatus{Directory: nil,
					Head:     "testing",
					Upstream: " ",
				}
			}(),
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		commandRunner := MockGitCommandRunner{statusResponse: tc.rawStatus}
		git := New(commandRunner)

		result := git.GetStatus()

		if result.Head != tc.expected.Head {
			t.Fatalf("Wrong head. Expected '%s' Got: '%s'", tc.expected.Head, result.Head)
		}

		if result.Upstream != tc.expected.Upstream {
			t.Fatalf("Wrong upstream. Expected '%s' Got: '%s'", tc.expected.Upstream, result.Upstream)
		}

		checkDir(result.Directory, tc.expected.Directory, t)
	}
}

func checkDir(directory, expectedDirectory *Directory, t *testing.T) {
	if expectedDirectory == nil {
		return
	}

	if directory.Name != expectedDirectory.Name {
		t.Fatalf("Wrong name name for directory. Expected: '%s' Got: '%s'", expectedDirectory.Name, directory.Name)
	}

	if len(directory.Directories) != len(expectedDirectory.Directories) {
		t.Fatalf("Wrong number of subdirectories for directory '%s'. Expected: '%d' Got: '%d'", directory.Name, len(expectedDirectory.Directories), len(directory.Directories))
	}

	if len(directory.Files) != len(expectedDirectory.Files) {
		t.Fatalf("Wrong number of files for directory '%s'. Expected: '%d' Got: '%d'", directory.Name, len(expectedDirectory.Files), len(directory.Files))
	}

	if directory.Parent != nil {
		if directory.Parent.Name != expectedDirectory.Parent.Name {
			t.Fatalf("Wrong parent name for directory '%s'. Expected: '%s' Got: '%s'", directory.Name, expectedDirectory.Parent.Name, directory.Parent.Name)
		}
	} else if directory.Parent == nil && expectedDirectory.Parent != nil {
		t.Fatalf("Expected no parent parent for directory '%s'.", directory.Name)
	}

	for i, expectedFile := range expectedDirectory.Files {
		file := directory.Files[i]
		checkFile(file, expectedFile, t)
	}

	for i, expectedSubDirectory := range expectedDirectory.Directories {
		subDirectory := directory.Directories[i]
		checkDir(subDirectory, expectedSubDirectory, t)
	}
}

func checkFile(file, expectedFile File, t *testing.T) {
	if file.Name != expectedFile.Name {
		t.Fatalf("Wrong name for file. Expected: '%s' Got: '%s'", expectedFile.Name, file.Name)
	}

	if file.Parent.Name != expectedFile.Parent.Name {
		t.Fatalf("Wrong parent name for file '%s'. Expected: '%s' Got: '%s'", file.Name, expectedFile.Parent.Name, file.Parent.Name)
	}

	if file.Dirpath != expectedFile.Dirpath {
		t.Fatalf("Wrong dirpath for file '%s'. Expected: '%s' Got: '%s'", file.Name, expectedFile.Dirpath, file.Dirpath)
	}

	if file.IndexStatus != expectedFile.IndexStatus {
		t.Fatalf("Wrong index status for file '%s'. Expected: '%d' Got: '%d'", file.Name, expectedFile.IndexStatus, file.IndexStatus)
	}

	if file.WorktreeStatus != expectedFile.WorktreeStatus {
		t.Fatalf("Wrong work tree status for file '%s'. Expected: '%d' Got: '%d'", file.Name, expectedFile.WorktreeStatus, file.WorktreeStatus)
	}
}
