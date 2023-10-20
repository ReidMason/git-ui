package ui

import (
	"fmt"
	"git-ui/internal/colours"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/state"

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
)

func RenderHeader(header string, viewWidth int) string {
	headerStyling := HeaderStyle.
		Width(viewWidth - HeaderStyle.GetHorizontalBorderSize())

	return headerStyling.Render(header)
}

func RenderFileTree(filetree filetree.FileTree, state state.State, width, height int) string {
	width = width - BorderStyle.GetHorizontalBorderSize()
	fileTreeString := lipgloss.NewStyle().
		MaxWidth(width).
		Render(filetree.Render())
	fileTreeString = lipgloss.NewStyle().Width(width).Height(height).Render(fileTreeString)

	style := BorderStyle
	if state.DiffsFocused() {
		style = BorderStyleInactive
	}

	return style.Render(fileTreeString)
}

func getColumnWidth(viewWidth int) int {
	return viewWidth / 12
}

func GetDiffDimensions(viewWidth, viewHeight int) (int, int) {
	headerHeight := 5
	footerHeight := 3
	return getColumnWidth(viewWidth) * 5, viewHeight - headerHeight - footerHeight
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

func renderDiffs(diffs string, state state.State) string {
	viewWidth := state.ViewWidth()
	viewHeight := state.ViewHeight()

	diffWidth, _ := GetDiffDimensions(viewWidth, viewHeight)
	diffWidth *= 2

	diffWidth -= BorderStyle.GetHorizontalBorderSize()
	diffs = lipgloss.NewStyle().MaxWidth(diffWidth).Render(diffs)
	diffs = lipgloss.NewStyle().Width(diffWidth).Render(diffs)

	style := BorderStyle
	if !state.DiffsFocused() {
		style = BorderStyleInactive
	}

	return style.Render(diffs)
}

func RenderMainView(fileTree filetree.FileTree, diffs string, statusbar string, state state.State) string {
	viewWidth := state.ViewWidth()
	viewHeight := state.ViewHeight()

	header := RenderHeader("Git-UI", viewWidth)

	diffs = renderDiffs(diffs, state)
	diffWidth, diffHeight := GetDiffDimensions(viewWidth, viewHeight)
	diffWidth *= 2

	leftoverWidth := viewWidth - diffWidth
	fileTreeString := RenderFileTree(fileTree, state, leftoverWidth, diffHeight)

	mainBody := lipgloss.JoinHorizontal(0, fileTreeString, diffs)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainBody, statusbar)
}
