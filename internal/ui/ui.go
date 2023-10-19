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

	HeaderStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Inherit(BorderStyle)
)

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
	return BorderStyle.Render(fileTreeString)
}

func getColumnWidth(viewWidth int) int {
	return viewWidth / 12
}

func GetDiffDimensions(viewWidth, viewHeight int) (int, int) {
	headerHeight := 5
	footerHeight := 3
	return getColumnWidth(viewWidth) * 5, viewHeight - headerHeight - footerHeight
}

func RenderStatusBar(state state.State, commitTextInput textinput.Model) string {
	status := state.GetGitStatus()
	viewWidth := state.GetViewWidth()

	ahead := lipgloss.NewStyle().Foreground(lipgloss.Color(colours.Green)).Render(fmt.Sprint(status.Ahead))
	behind := lipgloss.NewStyle().Foreground(lipgloss.Color(colours.Red)).Render(fmt.Sprint(status.Behind))
	output := fmt.Sprintf("%s | %s↑ | %s↓ ", status.Upstream, ahead, behind)

	width := viewWidth - BorderStyle.GetHorizontalBorderSize()
	output = lipgloss.PlaceHorizontal(width-50, lipgloss.Right, output)

	commitInput := commitTextInput.View()
	if !commitTextInput.Focused() {
		commitInput = ""
	}
	commitInput = lipgloss.NewStyle().Width(50).Render(commitInput)
	commitInput = lipgloss.NewStyle().MaxWidth(50).Render(commitInput)
	output = lipgloss.JoinHorizontal(0, commitInput, output)

	output = lipgloss.NewStyle().MaxWidth(width).Render(output)
	output = lipgloss.NewStyle().Width(width).Render(output)
	output = BorderStyle.Render(output)
	return output
}

func RenderMainView(viewWidth, viewHeight int, fileTree filetree.FileTree, diffs string, statusbar string) string {
	header := RenderHeader("Git-UI", viewWidth)

	diffWidth, diffHeight := GetDiffDimensions(viewWidth, viewHeight)
	diffWidth *= 2
	usedWidth := diffWidth

	diffWidth -= BorderStyle.GetHorizontalBorderSize()
	diffs = lipgloss.NewStyle().MaxWidth(diffWidth).Render(diffs)
	diffs = lipgloss.NewStyle().Width(diffWidth).Render(diffs)

	diffs = BorderStyle.Render(diffs)

	leftoverWidth := viewWidth - usedWidth
	fileTreeString := RenderFileTree(fileTree, leftoverWidth, diffHeight)

	mainBody := lipgloss.JoinHorizontal(0, fileTreeString, diffs)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainBody, statusbar)
}
