package ui

import (
	"git-ui/internal/colours"
	filetree "git-ui/internal/fileTree"

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
	return getColumnWidth(viewWidth) * 5, viewHeight - headerHeight
}

func RenderMainView(viewWidth, viewHeight int, fileTree filetree.FileTree, diffs string) string {
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

	return lipgloss.JoinVertical(lipgloss.Left, header, mainBody)
}
