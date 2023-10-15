package git

import (
	gitcommands "git-ui/internal/git_commands"
	"git-ui/internal/utils"
	"strings"
)

type GitInterface interface {
	GetStatus() *Directory
	GetDiff(filepath string) Diff
}

type Git struct {
	commandRunner gitcommands.GitCommandRunner
}

func New(commandRunner gitcommands.GitCommandRunner) Git {
	return Git{commandRunner: commandRunner}
}

type StatusChangeType rune

const (
	changed   StatusChangeType = '1'
	copied    StatusChangeType = '2'
	unmerged  StatusChangeType = 'u'
	untracked StatusChangeType = '?'
)

type GitStatus struct {
	Directory *Directory
}

func (g Git) GetStatus() GitStatus {
	rawStatus := g.commandRunner.GetStatus()
	lines := strings.Split(rawStatus, "\n")

	// First four lines are branch data so skip them for now
	lines = lines[4:]

	directory := newDirectory("Root", ".", nil)

	for _, line := range lines {
		firstRune, lineString := utils.TrimFirstRune(line)
		lineString = strings.TrimSpace(lineString)
		changeType := StatusChangeType(firstRune)

		if changeType == changed {
			file := parseChangedStatusLine(lineString)
			addFile(directory, strings.Split(file.Dirpath, "/"), make([]string, 0), file)
		}
		//   else if changeType == copied {
		//
		// } else if changeType == unmerged {
		//
		// } else if changeType == untracked {
		//
		// }
	}

	return GitStatus{Directory: directory}
}

func parseChangedStatusLine(line string) File {
	sections := strings.Split(line, " ")

	statusIndicators := []rune(sections[0])

	return newFile(sections[7], statusIndicators[0], statusIndicators[1])
}

func (g Git) Stage(filepath string)       { g.commandRunner.Stage(filepath) }
func (g Git) Unstage(filepath string)     { g.commandRunner.Unstage(filepath) }
func (g Git) Commit(commitMessage string) { g.commandRunner.Commit(commitMessage) }

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

func newDiff() Diff {
	return Diff{
		Diff1: make([]DiffLine, 0),
		Diff2: make([]DiffLine, 0),
	}
}

func DiffToString(difflines []DiffLine) string {
	content := ""
	for _, line := range difflines {
		content += line.Content + "\n"
	}

	return content
}

type DiffLine struct {
	Content string
	Type    DiffType
}

func newDiffLine(content string, diffType DiffType) DiffLine {
	return DiffLine{Content: content, Type: diffType}
}

func (g Git) GetDiff(filepath string) Diff {
	diffString := g.commandRunner.GetDiff(filepath)
	diff := newDiff()
	if len(diffString) == 0 {
		return diff
	}

	lines := strings.Split(diffString, "\n")[5:]

	changes := 0
	blankDiffLine := newDiffLine("", Blank)
	for _, line := range lines {
		indicator, lineString := utils.TrimFirstRune(line)
		if indicator == '-' {
			diffLine := newDiffLine(lineString, Removal)
			diff.Diff1 = append(diff.Diff1, diffLine)
			changes--
		} else if indicator == '+' {
			diffLine := newDiffLine(lineString, Addition)
			diff.Diff2 = append(diff.Diff2, diffLine)
			changes++
		} else {
			for changes < 0 {
				diff.Diff2 = append(diff.Diff2, blankDiffLine)
				changes++
			}

			for changes > 0 {
				diff.Diff1 = append(diff.Diff1, blankDiffLine)
				changes--
			}

			diffLine := newDiffLine(lineString, Neutral)
			diff.Diff1 = append(diff.Diff1, diffLine)
			diff.Diff2 = append(diff.Diff2, diffLine)
		}
	}

	for changes < 0 {
		diff.Diff2 = append(diff.Diff2, blankDiffLine)
		changes++
	}

	for changes > 0 {
		diff.Diff1 = append(diff.Diff1, blankDiffLine)
		changes--
	}

	return diff
}
