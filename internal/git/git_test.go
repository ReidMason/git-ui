package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type MockGitCommandRunner struct {
	diffRespponse string
}

func (g MockGitCommandRunner) Stage(filepath string)          {}
func (g MockGitCommandRunner) Unstage(filepath string)        {}
func (g MockGitCommandRunner) Commit(commitMessage string)    {}
func (g MockGitCommandRunner) GetDiff(filepath string) string { return g.diffRespponse }
func (g MockGitCommandRunner) GetStatus() string              { return "" }

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

		commandRunner := MockGitCommandRunner{diffRespponse: tc.rawDiff}
		git := New(commandRunner)
		result := git.GetDiff("path/to/file")

		for i, expectedDiffLine := range tc.expectedResult.Diff1 {
			resultDiffLine := result.Diff1[i]

			if !cmp.Equal(expectedDiffLine, resultDiffLine) {
				t.Fatalf("Diff1 failed line %d.\nExpected: '%v'\n     Got: '%v'", i+1, expectedDiffLine, resultDiffLine)
			}
		}

		for i, expectedDiffLine := range tc.expectedResult.Diff2 {
			resultDiffLine := result.Diff2[i]

			if expectedDiffLine.Content != resultDiffLine.Content {
				t.Fatalf("Diff2 failed line %d.\nExpected: '%s'\n     Got: '%s'", i+1, expectedDiffLine.Content, resultDiffLine.Content)
			}
		}
	}
}

//	func TestGetStatus(t *testing.T) {
//		rawStatus := `# branch.oid c86e7ed35f16570194c2308a2f8cb53155d0440d
//
// # branch.head main
// # branch.upstream origin/main
// # branch.ab +0 -0
// 1 .M N... 100644 100644 100644 51d742a142700c40e5d5d4915b44da5d238bef81 51d742a142700c40e5d5d4915b44da5d238bef81 internal/git/git.go
// 1 .M N... 100644 100644 100644 8508f049bcb61d4c52d92e5a4c9a71051f00bcba 8508f049bcb61d4c52d92e5a4c9a71051f00bcba internal/git/git_test.go
// 1 .M N... 100644 100644 100644 c789db6decaa4c7af3d5eb2214aea59f430dd5b1 c789db6decaa4c7af3d5eb2214aea59f430dd5b1 internal/utils/utils.go
// 1 M. N... 100644 100644 100644 1cdd739f6591c3aca07eab977748142a1ba14056 c345bc6f17650da4f51350e8faa56e4f4c61663e main.go
// ? internal/styling/styling.go`
//
//		result := *GetStatus(rawStatus)
//
//		expected := Directory{
//			Name:     "Root",
//			expanded: true,
//			Parent:   nil,
//			Files: []File{
//				{
//					Name:           "main.go",
//					Dirpath:        ".",
//					indexStatus:    77,
//					workTreeStatus: 46,
//				},
//			},
//			Directories: []*Directory{
//				{
//					Name:     "internal",
//					expanded: true,
//					Directories: []*Directory{
//						{
//							Name:        "git",
//							Directories: make([]*Directory, 0),
//							expanded:    true,
//							Files: []File{
//								{
//									Name:           "git.go",
//									Dirpath:        "internal/git",
//									indexStatus:    46,
//									workTreeStatus: 77,
//								},
//								{
//									Name:           "git_test.go",
//									Dirpath:        "internal/git",
//									indexStatus:    46,
//									workTreeStatus: 77,
//								},
//							},
//						},
//						{
//							Name:        "utils",
//							Directories: make([]*Directory, 0),
//							expanded:    true,
//							Files: []File{
//								{
//									Name:           "utils.go",
//									Dirpath:        "internal/utils",
//									indexStatus:    46,
//									workTreeStatus: 77,
//								},
//							},
//						},
//					},
//					Files: make([]File, 0),
//				},
//
//				// {
//				// 	Name:        "styling",
//				// 	Directories: nil,
//				// 	Files: []File{
//				// 		{
//				// 			Name:           "styling.go",
//				// 			Dirpath:        "internal/styling",
//				// 			IndexStatus:    46,
//				// 			WorkTreeStatus: 77,
//				// 		},
//				// 	},
//				// },
//			},
//		}
//
//		// checkDir(result, expected, t)
//
//		// if !cmp.Equal(result, expected) {
//		// 	t.Fatal("Wrong file path output")
//		// }
//	}
func checkDir(directory, expectedDirectory Directory, t *testing.T) {
	if !cmp.Equal(directory, expectedDirectory) {
		t.Fatalf("Directory mismatch. '%s' and '%s'", expectedDirectory.Name, directory.Name)
	}

	for i, expectedSubDirectory := range expectedDirectory.Directories {
		subDirectory := directory.Directories[i]
		checkDir(*subDirectory, *expectedSubDirectory, t)
	}
}
