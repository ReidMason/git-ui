package styling

import (
	"fmt"
	"git-ui/internal/git"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	ColumnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	HeaderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Align(lipgloss.Position(0.5))

	GreyOutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4b5161"))
)

func StyleLine(line git.DiffLine, width int) string {
	lineString := line.Content

	additionStyle := lipgloss.NewStyle().
		ColorWhitespace(true).
		Background(lipgloss.Color("#3f534f"))

	removalStyle := lipgloss.NewStyle().
		ColorWhitespace(true).
		Background(lipgloss.Color("#6f2e2d"))

	blankStyle := lipgloss.NewStyle().
		ColorWhitespace(true).
		Background(lipgloss.Color("#31343b"))

	if line.Type == git.Removal {
		lineString = removalStyle.Render(lipgloss.PlaceHorizontal(width, lipgloss.Left, lineString))
	} else if line.Type == git.Addition {
		lineString = additionStyle.Render(lipgloss.PlaceHorizontal(width, lipgloss.Left, lineString))
	} else if line.Type == git.Blank {
		lineString = blankStyle.Render(lipgloss.PlaceHorizontal(width, lipgloss.Left, lineString))
	}

	return lineString
}

func StyleDiff(diff []git.DiffLine, width int) string {
	diffString := ""
	count := 1

	numberOfLines := 0
	for _, line := range diff {
		if line.Type != git.Blank {
			numberOfLines++
		}
	}
	lineNumberPadding := len(fmt.Sprint(numberOfLines))

	for _, line := range diff {
		// TODO: This line number padding is probably really slow it could do with improving
		isBlank := line.Type == git.Blank
		lineNumber := strings.Repeat(" ", lineNumberPadding)
		if !isBlank {
			lengthOfCurrentNumber := len(fmt.Sprint(count))
			lineNumber = strings.Repeat(" ", lineNumberPadding-lengthOfCurrentNumber)
			lineNumber += fmt.Sprint(count)
			count++
		}
		lineNumber += "│"

		diffString += GreyOutStyle.Render(lineNumber) + StyleLine(line, width) + "\n"
	}

	return diffString
}
