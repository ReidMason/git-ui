package ui

import (
	"fmt"
	"git-ui/internal/colours"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/git"
	"git-ui/internal/state"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	BorderStyle = lipgloss.
			NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colours.Blue))

	BorderStyleInactive = BorderStyle.
				Copy().
				BorderForeground(lipgloss.Color(colours.Overlay0))

	HeaderStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Inherit(BorderStyle)

	DiffAddition = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colours.Base)).
			Background(lipgloss.Color(colours.Green)).
			Width(800)

	DiffRemoval = DiffAddition.Copy().
			Background(lipgloss.Color(colours.Red))

	DiffBlank = DiffAddition.Copy().
			Background(lipgloss.Color(colours.Surface0))

	GutterNumber = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colours.Surface2))

	LineSymbol = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colours.Overlay2))
)

func styleDiffLine(diffLine git.DiffLine) string {
	lineString := diffLine.Content
	if diffLine.Type == git.Addition {
		return DiffAddition.Render(lineString)
	} else if diffLine.Type == git.Removal {
		return DiffRemoval.Render(lineString)
	} else if diffLine.Type == git.Blank {
		return DiffBlank.Render(lineString)
	}

	return lineString
}

func buildLineSymbol(line git.DiffLine) string {
	lineSymbol := " "
	if line.Type == git.Addition {
		lineSymbol = "+"
	} else if line.Type == git.Removal {
		lineSymbol = "-"
	}

	return LineSymbol.Render(lineSymbol)
}

func DiffToString(difflines []git.DiffLine) string {
	gutterNumberPadding := 0
	for _, line := range difflines {
		if line.Type != git.Blank {
			gutterNumberPadding++
		}
	}
	gutterNumberPadding = len(fmt.Sprint(gutterNumberPadding))

	content := ""
	i := 1
	for _, line := range difflines {
		gutter := ""
		if line.Type != git.Blank {
			gutterNumber := fmt.Sprintf("%*d", gutterNumberPadding, i)
			gutter += GutterNumber.Render(gutterNumber)
			i++
		} else {
			gutter += strings.Repeat(" ", gutterNumberPadding)
		}

		content += gutter + buildLineSymbol(line) + GutterNumber.Render("│ ") + styleDiffLine(line) + "\n"
	}

	return content
}

func RenderHeader(header string, viewWidth int) string {
	headerStyling := HeaderStyle.
		Width(viewWidth - HeaderStyle.GetHorizontalBorderSize())

	return headerStyling.Render(header)
}

func RenderFileTree(filetree filetree.FileTree, width, height int) string {
	width = width - BorderStyle.GetHorizontalBorderSize()
	fileTreeString := lipgloss.NewStyle().
		MaxWidth(width).
		Render(filetree.Render())
	fileTreeString = lipgloss.NewStyle().Width(width).Height(height).Render(fileTreeString)

	style := BorderStyle
	if !filetree.Focused() {
		style = BorderStyleInactive
	}

	return style.Render(fileTreeString)
}

func GetDiffDimensions(viewWidth, viewHeight int) (int, int) {
	footerHeight := 3
	verticalBorderSize := BorderStyle.GetVerticalBorderSize()
	diffHeight := viewHeight - verticalBorderSize - footerHeight

	diffWidth := int(float32(viewWidth) * 0.4)

	return diffWidth, diffHeight
}

func GetFooterTextContent(state state.State) string {
	status := state.GetGitStatus()

	ahead := lipgloss.NewStyle().Foreground(lipgloss.Color(colours.Green)).Render(fmt.Sprint(status.Ahead))
	behind := lipgloss.NewStyle().Foreground(lipgloss.Color(colours.Red)).Render(fmt.Sprint(status.Behind))
	return fmt.Sprintf(" %s | %s | %s | %s ", status.Head, status.Upstream, ahead, behind)
}

func RenderStatusBar(state state.State, commitTextInput textinput.Model) string {
	width := state.ViewWidth() - BorderStyle.GetHorizontalBorderSize()

	output := GetFooterTextContent(state)
	outputLength := lipgloss.Width(output)

	commitInput := commitTextInput.View()
	if !commitTextInput.Focused() {
		commitInput = ""
	}

	commitInput = lipgloss.NewStyle().MaxWidth(width - outputLength).Render(commitInput)
	commitInput = lipgloss.NewStyle().Width(width - outputLength).Render(commitInput)

	commitInputLength := lipgloss.Width(commitInput)
	output = lipgloss.PlaceHorizontal(width-commitInputLength, lipgloss.Right, output)
	output = lipgloss.JoinHorizontal(0, commitInput, output)

	output = lipgloss.NewStyle().MaxWidth(width).Render(output)
	output = lipgloss.NewStyle().Width(width).Render(output)
	output = BorderStyle.Render(output)
	return output
}

func getDiffsBorderStyle(state state.State) lipgloss.Style {
	if state.DiffsFocused() {
		return BorderStyle
	}

	return BorderStyleInactive
}

func GetDiffsDividerStyle(state state.State) lipgloss.Style {
	if state.DiffsFocused() {
		return lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(colours.Blue))

	}

	return lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(colours.Overlay0))
}

func renderDiffs(diffs string, state state.State) string {
	viewWidth := state.ViewWidth()
	viewHeight := state.ViewHeight()

	diffWidth, _ := GetDiffDimensions(viewWidth, viewHeight)
	diffWidth *= 2

	diffWidth -= BorderStyle.GetHorizontalBorderSize()
	diffs = lipgloss.NewStyle().MaxWidth(diffWidth).Render(diffs)
	diffs = lipgloss.NewStyle().Width(diffWidth).Render(diffs)

	style := getDiffsBorderStyle(state)

	return style.Render(diffs)
}

func RenderMainView(fileTree filetree.FileTree, diffs string, statusbar string, state state.State) string {
	viewWidth := state.ViewWidth()
	viewHeight := state.ViewHeight()

	diffs = renderDiffs(diffs, state)
	diffWidth, diffHeight := GetDiffDimensions(viewWidth, viewHeight)
	diffWidth *= 2

	leftoverWidth := viewWidth - diffWidth
	fileTreeString := RenderFileTree(fileTree, leftoverWidth, diffHeight)

	mainBody := lipgloss.JoinHorizontal(0, fileTreeString, diffs)

	return lipgloss.JoinVertical(lipgloss.Left, mainBody, statusbar)
}
