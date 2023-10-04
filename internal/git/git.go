package git

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"unicode/utf8"
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

		letter, lineString := trimFirstRune(line)
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

func trimFirstRune(s string) (rune, string) {
	r, i := utf8.DecodeRuneInString(s)
	return r, s[i:]
}

func GetRawDiff(filepath string) string {
	cmd := exec.Command("git", "diff", "-U1000", filepath)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal("Failed to get git diff")
	}

	return strings.ReplaceAll(out.String(), "\t", "   ")
}

func GetStagedFiles(rawStaged string) []string {
	return strings.Split(rawStaged, "\n")
}

func GetRawStaged() string {
	cmd := exec.Command("git", "diff:", "--name-only", "--cached")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal("Failed to get staged files")
	}

	return strings.ReplaceAll(out.String(), "\t", "   ")
}
