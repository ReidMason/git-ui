package filetree

import (
	"git-ui/internal/git"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func StyleFileTreeLine(file FileTreeItem, selected, focused bool) lipgloss.Style {
	style := lipgloss.NewStyle()

	if file.IsFullyStaged() {
		style = style.Foreground(lipgloss.Color("#a6e3a1"))
	} else {
		style = style.Foreground(lipgloss.Color("#f38ba8"))
	}

	if selected && focused {
		return style.Copy().Background(lipgloss.Color("8"))
	} else if selected {
		return style.Copy().Background(lipgloss.Color("0"))
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
	GetFilePath() string
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
	currentLine   int
	IsFocused     bool
}

func New(directory git.Directory) FileTree {
	return FileTree{
		fileTreeLines: newFileTreeLines(directory, make([]FileTreeLine, 0), -1),
		currentLine:   0,
		IsFocused:     true,
	}
}

func (ft FileTree) Update(msg tea.Msg) (FileTree, tea.Cmd) {
	var cmd tea.Cmd
	ft, cmd = ft.updateAsModel(msg)
	return ft, cmd
}

func (ft FileTree) updateAsModel(msg tea.Msg) (FileTree, tea.Cmd) {
	var cmd tea.Cmd
	if !ft.IsFocused {
		return ft, cmd
	}

	downKeymap := key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down/j", "Down"),
	)

	upKeymap := key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up/k", "Up"),
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, downKeymap):
			ft.cursorDown()
		case key.Matches(msg, upKeymap):
			ft.cursorUp()
		}
	}

	return ft, cmd
}

func (ft FileTree) getDisplayedLines() int {
	count := -1
	for i := 1; i < len(ft.fileTreeLines); i++ {
		line := ft.fileTreeLines[i]
		if !line.Item.IsExpanded() {
			i += line.Item.Children()
		}

		count++
	}

	return count
}

func (ft *FileTree) cursorDown() {
	ft.currentLine = min(ft.currentLine+1, ft.getDisplayedLines())
}

func (ft *FileTree) cursorUp() {
	ft.currentLine = max(ft.currentLine-1, 0)
}

func (ft FileTree) GetIndex() int {
	return ft.currentLine
}

func (ft FileTree) GetSelectedFilepath() string {
	index := ft.currentLine
	for i := 1; i < len(ft.fileTreeLines); i++ {
		line := ft.fileTreeLines[i]

		if index == 0 {
			return line.Item.GetFilePath()
		}

		if !line.Item.IsExpanded() {
			i += line.Item.Children()
		}

		index--
	}

	return ""
}

func newFileTreeLines(directory git.Directory, fileTree []FileTreeLine, depth int) []FileTreeLine {
	newLine := newFileTreeLine(&directory, depth)
	fileTree = append(fileTree, newLine)

	for _, subDirectory := range directory.Directories {
		fileTree = newFileTreeLines(subDirectory, fileTree, depth+1)
	}

	for _, file := range directory.Files {
		newLine := newFileTreeLine(file, depth+1)
		fileTree = append(fileTree, newLine)
	}

	return fileTree
}

func (ft FileTree) Render() string {
	lines := ft.buildFileTreeString()
	output := ""

	for _, line := range lines {

		// if i == ft.currentLine {
		// 	output += ">"
		// } else {
		// 	output += " "
		// }
		output += line + "\n"
	}

	return output
}

func (ft FileTree) buildFileTreeString() []string {
	output := make([]string, 0)
	for i := 1; i < len(ft.fileTreeLines); i++ {
		line := ft.fileTreeLines[i]

		prefix := strings.Repeat("  ", line.Depth) + "-"
		lineString := prefix + line.Item.GetStatus() + " " + line.Item.GetName()
		selected := i-1 == ft.currentLine
		style := StyleFileTreeLine(line.Item, selected, ft.IsFocused)
		output = append(output, style.Render(lineString))

		if !line.Item.IsExpanded() {
			i += line.Item.Children()
		}
	}

	return output
}
