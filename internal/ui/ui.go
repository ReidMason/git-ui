package ui

import (
	"git-ui/internal/colours"
	filetree "git-ui/internal/fileTree"
	"log"

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

func RenderFileTree(filetree filetree.FileTree, width int) string {
	width = width - BorderStyle.GetHorizontalBorderSize()
	fileTreeString := lipgloss.NewStyle().
		MaxWidth(width).
		Render(filetree.Render())
	return BorderStyle.Render(fileTreeString)
}

func getColumnWidth(viewWidth int) int {
	return viewWidth / 12
}

func GetDiffDimensions(viewWidth, viewHeight int) (int, int) {
	return getColumnWidth(viewWidth) * 5, viewHeight - 5
}

func RenderMainView(viewWidth int, fileTree filetree.FileTree, diffs string) string {
	header := RenderHeader("Git-UI", viewWidth)

	diffWidth, _ := GetDiffDimensions(viewWidth, 0)
	diffWidth *= 2
	usedWidth := diffWidth

	diffWidth -= BorderStyle.GetHorizontalBorderSize()
	diffs = lipgloss.NewStyle().MaxWidth(diffWidth).Render(diffs)
	diffs = BorderStyle.Render(diffs)

	leftoverWidth := viewWidth - usedWidth
	log.Println("used width", usedWidth)
	log.Println("total width", viewWidth)
	log.Println("Leftover", leftoverWidth)
	fileTreeString := RenderFileTree(fileTree, leftoverWidth)

	mainBody := lipgloss.JoinHorizontal(0, fileTreeString, diffs)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainBody)
}
