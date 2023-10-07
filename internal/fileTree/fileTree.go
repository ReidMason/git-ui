package filetree

import (
	"git-ui/internal/git"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func StyleFileTreeLine(file FileTreeItem) lipgloss.Style {
	style := lipgloss.NewStyle()

	if file.IsFullyStaged() {
		style = style.Foreground(lipgloss.Color("#49b543"))
	} else {
		style = style.Foreground(lipgloss.Color("#f5b642"))
	}

	return style
}

type FileTreeItem interface {
	GetName() string
	IsExpanded() bool
	ToggleExpanded()
	Children() int
	GetStatus() string
	IsFullyStaged() bool
}

type FileTreeLine struct {
	Item  FileTreeItem
	Depth int
}

func newFileTreeLine(item FileTreeItem, depth int) FileTreeLine {
	return FileTreeLine{Item: item, Depth: depth}
}

type FileTree struct {
	fileTreeLines []FileTreeLine
}

func New(directory *git.Directory) FileTree {
	return FileTree{fileTreeLines: newFileTreeLines(directory, make([]FileTreeLine, 0), -1)}
}

func newFileTreeLines(directory *git.Directory, fileTree []FileTreeLine, depth int) []FileTreeLine {
	newLine := newFileTreeLine(directory, depth)
	fileTree = append(fileTree, newLine)

	for _, subDirectory := range directory.Directories {
		fileTree = newFileTreeLines(&subDirectory, fileTree, depth+1)
	}

	for _, file := range directory.Files {
		newLine := newFileTreeLine(file, depth+1)
		fileTree = append(fileTree, newLine)
	}

	return fileTree
}

func Render(lines []string) string {
	output := ""
	activeLine := 0

	for i, line := range lines {
		if i == activeLine {
			output += "> "
		} else {
			output += "  "
		}
		output += line + "\n"
	}

	return output
}

func (ft FileTree) BuildFileTreeString() []string {
	output := make([]string, 0)
	for i := 1; i < len(ft.fileTreeLines); i++ {
		line := ft.fileTreeLines[i]

		prefix := strings.Repeat(" ", line.Depth) + "-"
		lineString := prefix + line.Item.GetStatus() + " " + line.Item.GetName()
		style := StyleFileTreeLine(line.Item)
		output = append(output, style.Render(lineString))

		if !line.Item.IsExpanded() {
			i += line.Item.Children()
		}
	}

	return output
}
