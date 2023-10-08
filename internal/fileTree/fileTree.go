package filetree

import (
	"errors"
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
	IsVisible() bool
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
	cursorIndex   int
	IsFocused     bool
}

func New(directory *git.Directory) FileTree {
	return FileTree{
		fileTreeLines: newFileTreeLines(directory, make([]FileTreeLine, 0), -1)[1:],
		cursorIndex:   0,
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

	enterKeymap := key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Enter"),
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, downKeymap):
			ft.cursorDown()
		case key.Matches(msg, upKeymap):
			ft.cursorUp()
		case key.Matches(msg, enterKeymap):
			ft.handleEnter()
		}
	}

	return ft, cmd
}

func (ft *FileTree) handleEnter() {
	selectedLine, err := ft.getSelectedLine()
	if err != nil {
		return
	}

	switch lineItem := selectedLine.Item.(type) {
	case *git.Directory:
		lineItem.ToggleExpanded()
	}
}

func (ft *FileTree) cursorDown() {
	for i := ft.cursorIndex + 1; i < len(ft.fileTreeLines); i++ {
		if ft.updateCursorIndex(i) {
			return
		}
	}
}

func (ft *FileTree) cursorUp() {
	for i := ft.cursorIndex - 1; i >= 0; i-- {
		if ft.updateCursorIndex(i) {
			return
		}
	}
}

func (ft *FileTree) updateCursorIndex(newIndex int) bool {
	newSelectedLine := ft.fileTreeLines[newIndex]
	if newSelectedLine.Item.IsVisible() {
		ft.cursorIndex = newIndex
		return true
	}

	return false
}

func (ft FileTree) GetIndex() int {
	return ft.cursorIndex
}

func (ft FileTree) getSelectedLine() (FileTreeLine, error) {
	if len(ft.fileTreeLines) == 0 {
		var result FileTreeLine
		return result, errors.New("No file tree lines to display")
	}

	return ft.fileTreeLines[max(0, min(len(ft.fileTreeLines)-1, ft.cursorIndex))], nil
}

func (ft FileTree) GetSelectedFilepath() string {
	currentLine, err := ft.getSelectedLine()
	if err != nil {
		return ""
	}

	return currentLine.Item.GetFilePath()
}

func newFileTreeLines(directory *git.Directory, fileTree []FileTreeLine, depth int) []FileTreeLine {
	newLine := newFileTreeLine(directory, depth)
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
	return strings.Join(ft.buildFileTreeString(), "\n")
}

func (ft FileTree) buildFileTreeString() []string {
	output := make([]string, 0)

	if len(ft.fileTreeLines) == 0 {
		return append(output, "No changes")
	}

	for i := 0; i < len(ft.fileTreeLines); i++ {
		line := ft.fileTreeLines[i]
		if !line.Item.IsVisible() {
			continue
		}

		prefix := strings.Repeat("  ", line.Depth)

		icon := " "
		switch dir := line.Item.(type) {
		case *git.Directory:
			if dir.IsExpanded() {
				icon = "▼"
			} else {
				icon = "▶"
			}
		}

		prefix += icon
		lineString := prefix + line.Item.GetStatus() + " " + line.Item.GetName()

		selected := i == ft.cursorIndex
		style := StyleFileTreeLine(line.Item, selected, ft.IsFocused)
		output = append(output, style.Render(lineString))
	}

	return output
}
