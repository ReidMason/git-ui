package filetree

import (
	"git-ui/internal/styling"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileTreeElement interface {
	setSelected(selected bool)
	toggleExpanded()
	isVisible() bool
	getFilePath() string
	GetItem() FileTreeItem
}

type Directory struct {
	item        FileTreeItem
	parent      *Directory
	files       []*File
	directories []*Directory
	expanded    bool
	selected    bool
}

func newDirectory(parent *Directory, item FileTreeItem) Directory {
	return Directory{parent: parent, item: item, expanded: true}
}

func (d *Directory) setSelected(selected bool) { d.selected = selected }
func (d Directory) getFilePath() string        { return d.item.GetFilePath() }
func (d *Directory) toggleExpanded()           { d.expanded = !d.expanded }
func (d Directory) GetItem() FileTreeItem      { return d.item }
func (d Directory) isVisible() bool {
	if d.parent == nil {
		return true
	}

	return d.parent.isVisible() && d.parent.expanded
}

type File struct {
	parent   *Directory
	item     FileTreeItem
	selected bool
}

func (f *File) setSelected(selected bool) { f.selected = selected }
func (f File) getFilePath() string        { return f.item.GetFilePath() }
func (f File) GetItem() FileTreeItem      { return f.item }
func (f *File) toggleExpanded()           {}
func (f File) isVisible() bool            { return f.parent.isVisible() && f.parent.expanded }

type FileTreeLine struct {
	item       FileTreeItem
	depth      int
	isExpanded bool
}

type FileTreeItem interface {
	GetDisplay() string
	IsDirectory() bool
	GetFilePath() string
	GetDirectories() []FileTreeItem
	GetFiles() []FileTreeItem
}

type FileTree struct {
	fileTreeItems []FileTreeElement
	root          Directory
	cursorIndex   int
	width         int
	isFocused     bool
}

func New(directory FileTreeItem) FileTree {
	fileTree := FileTree{
		isFocused: true,
	}

	return fileTree.UpdateDirectoryTree(directory, "")
}

func (ft FileTree) UpdateDirectoryTree(directory FileTreeItem, selectedFilepath string) FileTree {
	cursorIndex := ft.buildTree(directory, selectedFilepath)
	ft.setCursorIndex(cursorIndex)
	return ft
}

func (ft FileTree) Update(msg tea.Msg, spaceCmd tea.Cmd) (FileTree, tea.Cmd) {
	if !ft.isFocused {
		return ft, nil
	}

	keyDown := key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down/j", "Down"),
	)

	keyUp := key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up/k", "Up"),
	)

	keyEnter := key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Enter"),
	)

	keySpace := key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "Space"),
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyDown):
			ft.handleKeyDown()
		case key.Matches(msg, keyUp):
			ft.handleKeyUp()
		case key.Matches(msg, keyEnter):
			ft.handleKeyEnter()
		case key.Matches(msg, keySpace):
			return ft, spaceCmd
		}
	case tea.WindowSizeMsg:
		ft.width = msg.Width
	}

	return ft, nil
}

func (ft *FileTree) handleKeyDown() {
	for i := ft.cursorIndex + 1; i < len(ft.fileTreeItems); i++ {
		if ft.fileTreeItems[i].isVisible() {
			ft.setCursorIndex(i)
			break
		}
	}
}

func (ft *FileTree) handleKeyUp() {
	for i := ft.cursorIndex - 1; i >= 0; i-- {
		if ft.fileTreeItems[i].isVisible() {
			ft.setCursorIndex(i)
			break
		}
	}
}

func (ft *FileTree) setCursorIndex(cursorIndex int) {
	if len(ft.fileTreeItems) == 0 {
		return
	}

	for _, item := range ft.fileTreeItems {
		item.setSelected(false)
	}

	ft.cursorIndex = cursorIndex
	ft.fileTreeItems[cursorIndex].setSelected(true)
}

func (ft *FileTree) handleKeyEnter() {
	selected := ft.fileTreeItems[ft.cursorIndex]
	selected.toggleExpanded()
}

func newFileTreeLine(item FileTreeItem, depth int) FileTreeLine {
	return FileTreeLine{item: item, depth: depth, isExpanded: true}
}

func (ft FileTree) Render() string {
	return strings.Join(ft.buildFileTreeString(), "\n")
}

func (ft *FileTree) buildTree(directory FileTreeItem, selectedFilepath string) int {
	selectedIndex := 0
	ft.root = newDirectory(nil, directory)
	ft.fileTreeItems = nil

	for _, subDirectory := range directory.GetDirectories() {
		res := ft.buildTreeR(&ft.root, subDirectory, selectedFilepath)
		if res > 0 {
			selectedIndex = res
		}
	}

	for _, file := range directory.GetFiles() {
		if file.GetFilePath() == selectedFilepath {
			selectedIndex = len(ft.fileTreeItems)

		}

		newFile := File{parent: &ft.root, item: file}
		ft.fileTreeItems = append(ft.fileTreeItems, &newFile)
		ft.root.files = append(ft.root.files, &newFile)
	}

	return selectedIndex
}

func (ft *FileTree) buildTreeR(parent *Directory, directory FileTreeItem, selectedFilepath string) int {
	selectedIndex := 0
	if directory.GetFilePath() == selectedFilepath {
		selectedIndex = len(ft.fileTreeItems)
	}

	newDirectory := newDirectory(parent, directory)
	ft.fileTreeItems = append(ft.fileTreeItems, &newDirectory)

	for _, subDirectory := range directory.GetDirectories() {
		res := ft.buildTreeR(&newDirectory, subDirectory, selectedFilepath)

		if res > 0 {
			selectedIndex = res
		}
	}

	for _, file := range directory.GetFiles() {
		if file.GetFilePath() == selectedFilepath {
			selectedIndex = len(ft.fileTreeItems)
		}

		newFile := File{parent: &newDirectory, item: file}
		ft.fileTreeItems = append(ft.fileTreeItems, &newFile)
		newDirectory.files = append(newDirectory.files, &newFile)
	}

	parent.directories = append(parent.directories, &newDirectory)

	return selectedIndex
}

func getIcon(directory Directory) string {
	if directory.expanded {
		return "▼"
	}

	return "▶"
}

func buildFileOutputString(file File, output []string, depth int, width int) []string {
	line := strings.Repeat(" ", depth+3)
	line += file.item.GetDisplay()

	if file.selected {
		line = getSelectedStyle(line, width)
	}

	return append(output, line)
}

func buildFileTreeElementOutputString(directory Directory, output []string, depth int, width int) []string {
	line := strings.Repeat("  ", depth) + getIcon(directory)
	line += " " + directory.item.GetDisplay()

	if directory.selected {
		line = getSelectedStyle(line, width)
	}

	output = append(output, line)

	if !directory.expanded {
		return output
	}

	for _, subDirectory := range directory.directories {
		output = buildFileTreeElementOutputString(*subDirectory, output, depth+1, width)
	}

	for _, file := range directory.files {
		output = buildFileOutputString(*file, output, depth+1, width)
	}

	return output
}

func (ft FileTree) GetSelectedFilepath() string {
	if len(ft.fileTreeItems) == 0 {
		return ""
	}

	selectedItem := ft.GetSelectedItem()
	return selectedItem.getFilePath()
}

func (ft FileTree) GetSelectedItem() FileTreeElement {
	return ft.fileTreeItems[ft.cursorIndex]
}

func (ft FileTree) buildFileTreeString() []string {
	output := make([]string, 0)

	if len(ft.root.files)+len(ft.root.directories) == 0 {
		return append(output, "No changes")
	}

	for _, subDirectory := range ft.root.directories {
		output = buildFileTreeElementOutputString(*subDirectory, output, 0, ft.width)
	}

	for _, file := range ft.root.files {
		output = buildFileOutputString(*file, output, -3, ft.width)
	}

	return output
}

func (ft *FileTree) SetFocused(focused bool) {
	ft.isFocused = focused
}

func getSelectedStyle(line string, width int) string {
	selectedStyling := lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color("8"))
	line = styling.TrimColourResetChar(line)
	line = lipgloss.PlaceHorizontal(width-2, lipgloss.Left, line)
	return selectedStyling.Render(line)

	//  else if selected {
	// 	return lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color("0"))
	// }
}
