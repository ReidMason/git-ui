package git

import (
	"strconv"
	"strings"

	gitcommands "github.com/ReidMason/git-ui/internal/git_commands"
	"github.com/ReidMason/git-ui/internal/utils"
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
	Head      string
	Upstream  string
	Ahead     int
	Behind    int
}

func parseHead(line string) string {
	return strings.TrimPrefix(line, "# branch.head ")
}

func parseAheadAndBehind(line string) (int, int) {
	elements := strings.Split(line, " ")
	ahead := 0
	behind := 0

	if len(elements) >= 3 {
		ahead, _ = strconv.Atoi(elements[2])
	}

	if len(elements) >= 4 {
		behind, _ = strconv.Atoi(elements[3])
	}

	return ahead, behind
}

func parseUpstream(line string) string {
	return strings.TrimPrefix(line, "# branch.upstream ")
}

func addStatusMetadata(line string, gitStatus GitStatus) GitStatus {
	if strings.HasPrefix(line, "# branch.head") {
		gitStatus.Head = parseHead(line)
		return gitStatus
	}

	if strings.HasPrefix(line, "# branch.upstream") {
		gitStatus.Upstream = parseUpstream(line)
		return gitStatus
	}

	if strings.HasPrefix(line, "# branch.ab") {
		gitStatus.Ahead, gitStatus.Behind = parseAheadAndBehind(line)
		return gitStatus
	}

	return gitStatus
}

func (g Git) GetStatus() GitStatus {
	rawStatus, err := g.commandRunner.GetStatus()
	if err != nil {
		return GitStatus{}
	}
	lines := strings.Split(rawStatus, "\n")

	gitStatus := GitStatus{
		Directory: newDirectory("Root", ".", nil),
		Upstream:  " ",
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			gitStatus = addStatusMetadata(line, gitStatus)
			continue
		}

		firstRune, lineString := utils.TrimFirstRune(line)
		lineString = strings.TrimSpace(lineString)
		changeType := StatusChangeType(firstRune)

		if changeType == changed {
			file := parseChangedStatusLine(lineString)
			addFile(gitStatus.Directory, strings.Split(file.Dirpath, "/"), make([]string, 0), file)
		} else if changeType == untracked {
			file := parseUntrackedStatusLine(lineString)
			addFile(gitStatus.Directory, strings.Split(file.Dirpath, "/"), make([]string, 0), file)
		} else if changeType == copied {
			file := parseCopiedStatusLine(lineString)
			addFile(gitStatus.Directory, strings.Split(file.Dirpath, "/"), make([]string, 0), file)
		}
	}

	gitStatus.Directory.Sort()

	return gitStatus
}

func parseCopiedStatusLine(line string) File {
	sections := strings.Split(line, " ")
	statusIndicators := []rune(sections[0])

	secondName := "?"
	if len(sections) >= 12 {
		secondName = sections[11]
	}

	return newFile(sections[8], statusIndicators[0], statusIndicators[1], secondName)
}

func parseChangedStatusLine(line string) File {
	sections := strings.Split(line, " ")

	statusIndicators := []rune(sections[0])

	return newFile(sections[7], statusIndicators[0], statusIndicators[1], "")
}

func parseUntrackedStatusLine(line string) File {
	filepath := strings.TrimPrefix(line, "? ")
	return newFile(filepath, '?', '?', "")
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

	lines := strings.Split(diffString, "\n")

	start := false
	changes := 0
	blankDiffLine := newDiffLine("", Blank)
	for _, line := range lines {
		if !start && !strings.HasPrefix(line, "@@") {
			continue
		} else if strings.HasPrefix(line, "@@") {
			start = true
			continue
		}

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
