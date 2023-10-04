package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestGetRawDiff(t *testing.T) {
// 	result := GetRawDiff()
//
// 	expected := `diff --git git-ui/testfile.txt git-ui/testfile.txt
// index 35b5809..4492ac6 100644
// --- git-ui/testfile.txt
// +++ git-ui/testfile.txt
// @@ -1,4 +1,5 @@
//  This is a test file
// -
//  These lines are committed now
// -I have added some more content
// +I have added this is a change more content
// +
// +This is a new thing`
//
// 	resultLines := strings.Split(result, "\n")
// 	for i, line := range strings.Split(expected, "\n") {
// 		resLine := resultLines[i]
//
// 		if line != resLine {
// 			t.Fatalf("Line match failed. Expected: '%s' Got: '%s'", line, resLine)
// 		}
// 	}
//
// 	if len(result) != len(expected) {
// 		t.Fatalf("Length difference. Expected: %d Got: %d", len(expected), len(result))
// 	}
// }

func TestGetDiff(t *testing.T) {
	rawDiff := `diff --git a/git-ui/testfile.txt b/git-ui/testfile.txt
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
 }`

	result := GetDiff(rawDiff)

	expected := Diff{
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
	}

	t.Log("----- Start diff1")
	for i, expectedDiffLine := range expected.Diff1 {
		resultDiffLine := result.Diff1[i]
		t.Log(resultDiffLine.Content)

		if !cmp.Equal(expectedDiffLine, resultDiffLine) {
			t.Fatalf("Diff1 failed line %d.\nExpected: '%v'\n     Got: '%v'", i+1, expectedDiffLine, resultDiffLine)
		}
	}

	t.Log("----- Start diff2")
	for i, expectedDiffLine := range expected.Diff2 {
		resultDiffLine := result.Diff2[i]
		t.Log(resultDiffLine.Content)

		if expectedDiffLine.Content != resultDiffLine.Content {
			t.Fatalf("Diff2 failed line %d.\nExpected: '%s'\n     Got: '%s'", i+1, expectedDiffLine.Content, resultDiffLine.Content)
		}
	}

	// if !reflect.DeepEqual(result, expected) {
	// 	t.Fatal("Expected diff length doesn't match")
	// }
}

func TestGetFiles(t *testing.T) {
	rawFiles := `internal/git/git.go
main.go
internal/git/git_test.go
internal/utils/utils.go`

	result := GetFiles(rawFiles, "")

	expected := []File{
		{
			Name:      "internal",
			Status:    FileStatus(None),
			Directory: true,
			Files: []File{
				{
					Name:      "git",
					Status:    FileStatus(None),
					Directory: true,
					Files: []File{
						{
							Name:      "git.go",
							Status:    FileStatus(Staged),
							Directory: false,
						},
						{
							Name:      "git_test.go",
							Status:    FileStatus(Staged),
							Directory: false,
						},
					},
				},
				{
					Name:      "utils",
					Status:    FileStatus(None),
					Directory: true,
					Files: []File{
						{
							Name:      "utils.go",
							Status:    FileStatus(Staged),
							Directory: false,
						},
					},
				},
			},
		},
		{
			Name:      "main.go",
			Status:    FileStatus(Staged),
			Directory: false,
		},
	}

	// s, _ := json.MarshalIndent(result, "", "\t")
	// t.Logf(string(s))

	f1 := expected[0]
	f2 := result[0]
	if !compareFile(f1, f2) {
		t.Fatalf("File comparison for 'internal' directory failed. Expected: '%v' Got: '%v'", f1, f2)
	}

	f1 = expected[0].Files[0]
	f2 = result[0].Files[0]
	if !compareFile(f1, f2) {
		t.Fatalf("File comparison for 'git' directory failed. Expected: '%v' Got: '%v'", f1, f2)
	}

	f1 = expected[0].Files[0].Files[0]
	f2 = result[0].Files[0].Files[0]
	if !compareFile(f1, f2) {
		t.Fatalf("File comparison for 'git.go' directory failed. Expected: '%v' Got: '%v'", f1, f2)
	}

	f1 = expected[0].Files[0].Files[1]
	f2 = result[0].Files[0].Files[1]
	if !compareFile(f1, f2) {
		t.Fatalf("File comparison for 'git_test.go' directory failed. Expected: '%v' Got: '%v'", f1, f2)
	}

	f1 = expected[0].Files[1]
	f2 = result[0].Files[1]
	if !compareFile(f1, f2) {
		t.Fatalf("File comparison for 'utils' directory failed. Expected: '%v' Got: '%v'", f1, f2)
	}

	f1 = expected[0].Files[1].Files[0]
	f2 = result[0].Files[1].Files[0]
	if !compareFile(f1, f2) {
		t.Fatalf("File comparison for 'utils.go' failed. Expected: '%v' Got: '%v'", f1, f2)
	}

	f1 = expected[1]
	f2 = result[1]
	if !compareFile(f1, f2) {
		t.Fatalf("File comparison for 'main.go' failed. Expected: '%v' Got: '%v'", f1, f2)
	}
}

func compareFile(file1, file2 File) bool {
	if file1.Name != file2.Name {
		return false
	}

	if file1.Directory != file2.Directory {
		return false
	}

	if file1.Status != file2.Status {
		return false
	}

	return true
}
