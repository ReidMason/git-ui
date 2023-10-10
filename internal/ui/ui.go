package ui

import (
	filetree "git-ui/internal/fileTree"

	"github.com/charmbracelet/lipgloss"
)

var (
	BorderStyle = lipgloss.
			NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	HeaderStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Inherit(BorderStyle)
)

func RenderHeader(header string, viewWidth int) string {
	headerStyling := HeaderStyle.
		Width(viewWidth - HeaderStyle.GetHorizontalBorderSize())
	return headerStyling.Render(header)
}

func RenderFileTree(filetree filetree.FileTree) string {
	return filetree.Render()
}

func RenderMainView(viewWidth int, fileTree filetree.FileTree) string {
	header := RenderHeader("Git-UI", viewWidth)
	fileTreeString := RenderFileTree(fileTree)

	return lipgloss.JoinVertical(lipgloss.Left, header, fileTreeString)
}
