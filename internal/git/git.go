package git

import (
	"git-ui/internal/utils"
	"log"
	"strings"
)

type DiffType int8

const (
	Removal DiffType = iota
	Addition
	Neutral
	Blank
)

type Diff struct {
	Diff1 []DiffLine
	Diff2 []DiffLine
}

type DiffLine struct {
	Content string
	Type    DiffType
}

func GetDiff(diffString string) Diff {
	lines := strings.Split(diffString, "\n")

	diff := Diff{Diff1: make([]DiffLine, 0), Diff2: make([]DiffLine, 0)}

	start := false
	removals := 0
	additions := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "@@") && strings.HasSuffix(line, "@@") {
			start = true
			continue
		}

		if !start {
			continue
		}

		letter, lineString := utils.TrimFirstRune(line)
		if letter == '-' {
			diffLine := DiffLine{Content: lineString, Type: Removal}
			diff.Diff1 = append(diff.Diff1, diffLine)
			removals++
		} else if letter == '+' {
			diffLine := DiffLine{Content: lineString, Type: Addition}
			diff.Diff2 = append(diff.Diff2, diffLine)
			if removals > 0 {
				removals--
			} else {
				additions++
			}
		} else {
			diffLine := DiffLine{Content: "", Type: Blank}
			for i := 0; i < removals; i++ {
				diff.Diff2 = append(diff.Diff2, diffLine)
			}
			removals = 0

			diffLine = DiffLine{Content: "", Type: Blank}
			for i := 0; i < additions; i++ {
				diff.Diff1 = append(diff.Diff1, diffLine)
			}
			additions = 0

			diffLine = DiffLine{Content: lineString, Type: Neutral}
			diff.Diff1 = append(diff.Diff1, diffLine)
			diff.Diff2 = append(diff.Diff2, diffLine)
		}
	}

	diffLine := DiffLine{Content: "", Type: Blank}
	for i := 0; i < removals; i++ {
		diff.Diff2 = append(diff.Diff2, diffLine)
	}

	diffLine = DiffLine{Content: "", Type: Blank}
	for i := 0; i < additions; i++ {
		diff.Diff1 = append(diff.Diff1, diffLine)
	}

	return diff
}

type File struct {
	Name  string
	Files []File
}

func newFile(filename string) File {
	return File{Name: filename, Files: make([]File, 0)}
}

func GetFiles(rawStagedFiles, rawUnstagedFiles string) []File {
	files := make([]File, 0)

	stagedFilepaths := strings.Split(rawStagedFiles, "\n")
	for _, filepath := range stagedFilepaths {
		files = addFile(files, filepath)
	}

	return files
}

func addFile(files []File, filepath string) []File {
	if !strings.Contains(filepath, "/") {
		log.Println("here at the file")
		files = append(files, newFile(filepath))
		return files
	}

	splitFilepath := strings.SplitN(filepath, "/", 2)
	filename := splitFilepath[0]
	filepath = splitFilepath[1]

	added := false
	for i, file := range files {
		if file.Name == filename {
			files[i].Files = addFile(file.Files, filepath)
			added = true
			break
		}
	}

	if added == false {
		parent := newFile(filename)
		parent.Files = addFile(parent.Files, filepath)
		files = append(files, parent)
	}

	return files
}

func GetRawDiff(filepath string) string {
	result, err := utils.RunCommand("git", "diff", "-U1000", filepath)

	if err != nil {
		log.Fatal("Failed to get git diff")
	}

	return result
}

func GetRawStaged() string {
	result, err := utils.RunCommand("git", "diff:", "--name-only", "--cached")
	if err != nil {
		log.Fatal("Failed to get staged files")
	}

	return result
}

func GetRawUnstaged() string {
	result, err := utils.RunCommand("git", "diff:", "--name-only")
	if err != nil {
		log.Fatal("Failed to get unstaged files")
	}

	return result
}
