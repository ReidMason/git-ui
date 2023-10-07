package git

import (
	"git-ui/internal/utils"
	"log"
	"path/filepath"
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

type FileStatus int8

const (
	Staged DiffType = iota
	Unstaged
	None
)

const (
	Changed   = '1'
	Copied    = '2'
	Unmerged  = 'u'
	Untracked = '?'
)

type Directory struct {
	Name        string
	Files       []File
	Directories []Directory
	expanded    bool
}

func newDirectory(name string) Directory {
	return Directory{Name: name, Files: make([]File, 0), Directories: make([]Directory, 0), expanded: true}
}

func (d Directory) GetName() string  { return d.Name }
func (d *Directory) ToggleExpanded() { d.expanded = !d.expanded }
func (d Directory) Children() int {
	count := len(d.Files)
	for _, directory := range d.Directories {
		count++
		count += directory.Children()
	}

	return count
}
func (d Directory) GetStatus() string { return "" }

func (d Directory) IsExpanded() bool {
	return d.expanded
}

func (d Directory) IsFullyStaged() bool {
	for _, file := range d.Files {
		if !file.IsFullyStaged() {
			return false
		}
	}

	for _, directory := range d.Directories {
		if !directory.IsFullyStaged() {
			return false
		}
	}

	return true
}

type File struct {
	Name           string
	Dirpath        string
	indexStatus    rune
	workTreeStatus rune
}

func (f File) GetName() string {
	return f.Name
}

func (f File) ToggleExpanded()     {}
func (f File) IsExpanded() bool    { return true }
func (f File) Children() int       { return 0 }
func (f File) IsFullyStaged() bool { return f.workTreeStatus == '.' }
func (f File) GetStatus() string   { return string(f.indexStatus) + string(f.workTreeStatus) }

func newFile(filePath string, indexStatus, workTreeStatus rune) File {
	dirpath, filename := filepath.Split(filePath)
	dirpath = filepath.Clean(dirpath)

	return File{Name: filename, Dirpath: dirpath, indexStatus: indexStatus, workTreeStatus: workTreeStatus}
}

func GetStatus(rawStatus string) Directory {
	lines := strings.Split(rawStatus, "\n")

	// First four lines are branch data so skip them for now
	lines = lines[4:]

	directory := newDirectory("Root")

	for _, line := range lines {
		changeType, line := utils.TrimFirstRune(line)
		line = strings.TrimSpace(line)

		if changeType == Changed {
			file := parseChangedLine(line)
			directory = addFile(directory, strings.Split(file.Dirpath, "/"), file)
		}
		//   else if changeType == Copied {
		//
		// } else if changeType == Unmerged {
		//
		// } else if changeType == Untracked {
		//
		// }
	}

	return directory
}

func parseChangedLine(line string) File {
	sections := strings.Split(line, " ")

	statusIndicators := []rune(sections[0])

	return newFile(sections[7], statusIndicators[0], statusIndicators[1])
}

func addFile(directory Directory, dirpath []string, newFile File) Directory {
	if len(dirpath) == 0 || dirpath[0] == "." {
		directory.Files = append(directory.Files, newFile)
		return directory
	}

	for i, subdir := range directory.Directories {
		if subdir.Name == dirpath[0] {
			directory.Directories[i] = addFile(subdir, dirpath[1:], newFile)
			return directory
		}
	}

	newDir := newDirectory(dirpath[0])
	newDir = addFile(newDir, dirpath[1:], newFile)
	directory.Directories = append(directory.Directories, newDir)

	return directory
}

func GetRawDiff(filepath string) string {
	result, err := utils.RunCommand("git", "diff", "-U1000", filepath)

	if err != nil {
		log.Fatal("Failed to get git diff")
	}

	return result
}

func GetRawStatus() string {
	result, err := utils.RunCommand("git", "status", "-u", "--porcelain=v2", "--branch")
	if err != nil {
		log.Fatal("Failed to get git status")
	}

	return result
}
