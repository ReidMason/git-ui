package filetree

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Directory struct {
	item        FileTreeItem
	parent      *Directory
	files       []File
	directories []*Directory
	isExpanded  bool
	selected    bool
}

func newDirectory(parent *Directory, item FileTreeItem) Directory {
	return Directory{parent: nil, item: item, isExpanded: true}
}

type File struct {
	parent   *Directory
	item     FileTreeItem
	selected bool
}

type FileTreeLine struct {
	item       FileTreeItem
	depth      int
	isExpanded bool
}

// func (l FileTreeLine) IsVisible() bool {
// 	if l.Parent == nil {
// 		return true
// 	}
//
// 	return l.Parent.IsVisible() && d.Parent.Expanded
// }

type FileTreeItem interface {
	GetName() string
	Children() int
	IsDirectory() bool
	GetFilePath() string
	GetDirectories() []FileTreeItem
	GetFiles() []FileTreeItem
}

type FileTree struct {
	// fileTreeLines []FileTreeLine
	root        Directory
	cursorIndex int
	isFocused   bool
}

func New(directory FileTreeItem) FileTree {

	return FileTree{
		root: buildTree(directory),
		// fileTreeLines: buildFileTreeLines(rootDirectory),
		cursorIndex: 0,
		isFocused:   true,
	}
}

func (ft FileTree) Update(msg tea.Msg) {
	keyDown := key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down/j", "Down"),
	)

	keyUp := key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up/k", "Up"),
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyDown):
			ft.handleKeyDown()
		case key.Matches(msg, keyUp):
			ft.handleKeyUp()
		}
	}
}

func (ft FileTree) handleKeyDown() {
	log.Println("Down")
}

func (ft FileTree) handleKeyUp() {
	log.Println("Up")
}

func newFileTreeLine(item FileTreeItem, depth int) FileTreeLine {
	return FileTreeLine{item: item, depth: depth, isExpanded: true}
}

func (ft FileTree) Render() string {
	return strings.Join(ft.buildFileTreeString(), "\n")
}

func buildTree(directory FileTreeItem) Directory {
	rootDirectory := newDirectory(nil, directory)

	for _, subDirectory := range directory.GetDirectories() {
		buildTreeR(&rootDirectory, subDirectory)
	}

	for _, file := range directory.GetFiles() {
		newFile := File{parent: &rootDirectory, item: file}
		rootDirectory.files = append(rootDirectory.files, newFile)
	}

	return rootDirectory
}

func buildTreeR(parent *Directory, directory FileTreeItem) {
	newDirectory := newDirectory(parent, directory)

	for _, subDirectory := range directory.GetDirectories() {
		buildTreeR(&newDirectory, subDirectory)
	}

	for _, file := range directory.GetFiles() {
		newFile := File{parent: &newDirectory, item: file}
		newDirectory.files = append(newDirectory.files, newFile)
	}

	parent.directories = append(parent.directories, &newDirectory)
}
func getIcon(directory Directory) string {
	if directory.isExpanded {
		return "▼"
	}

	return "▶"
}

func buildFileOutputString(file File, output []string, depth int) []string {
	line := strings.Repeat("  ", depth+1) + file.item.GetName()

	if file.selected {
		selectedStyling := styleFileSelected(file.selected)
		line = selectedStyling.Render(line)
	}
	return append(output, line)
}

func buildFileTreeElementOutputString(directory Directory, output []string, depth int) []string {
	line := strings.Repeat("  ", depth) + getIcon(directory)
	line += " " + directory.item.GetName()

	if directory.selected {
		selectedStyling := styleFileSelected(directory.selected)
		line = selectedStyling.Render(line)
	}

	output = append(output, line)

	for _, subDirectory := range directory.directories {
		output = buildFileTreeElementOutputString(*subDirectory, output, depth+1)
	}

	for _, file := range directory.files {
		output = buildFileOutputString(file, output, depth+1)
	}

	return output
}

func (ft FileTree) buildFileTreeString() []string {
	output := make([]string, 0)

	if len(ft.root.files)+len(ft.root.directories) == 0 {
		return append(output, "No changes")
	}

	// Remove this it's only for testing
	ft.root.directories[0].selected = true

	return buildFileTreeElementOutputString(ft.root, output, 0)

	// for i := 0; i < len(ft.fileTreeLines); i++ {
	// 	line := ft.fileTreeLines[i]
	// 	// if !line.isVisible() {
	// 	// 	continue
	// 	// }
	//
	// 	prefix := strings.Repeat("  ", line.depth) + getIcon(line)
	//
	// 	lineString := line.item.GetName()
	// 	selected := i == ft.cursorIndex
	// 	if selected {
	// 		lineString = lipgloss.PlaceHorizontal(50, lipgloss.Left, lineString)
	// 	}
	//
	// 	// style := styleFileTreeLine(line.Item)
	// 	// lineString = prefix + style.Render(lineString)
	//
	// 	lineString = prefix + line.item.GetName()
	//
	// 	// if selected {
	// 	// 	selectedStyling := styleFileSelected(selected, ft.IsFocused)
	// 	// 	lineString = selectedStyling.Render(lineString)
	// 	// }
	//
	// 	output = append(output, lineString)
	// }
}

func styleFileSelected(selected bool) lipgloss.Style {
	if selected {
		return lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color("8"))
	}
	//  else if selected {
	// 	return lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color("0"))
	// }

	return lipgloss.NewStyle()
}

// func buildFileTreeLines(directory Directory) []FileTreeLine {
// 	return newFileTreeLines(directory, make([]FileTreeLine, 0), -1)[1:]
// }
//
// func newFileTreeLines(directory Directory, fileTree []FileTreeLine, depth int) []FileTreeLine {
// 	newLine := newFileTreeLine(directory, depth)
// 	fileTree = append(fileTree, newLine)
//
// 	for _, subDirectory := range directory.GetDirectories() {
// 		fileTree = newFileTreeLines(subDirectory, fileTree, depth+1)
// 	}
//
// 	for _, file := range directory.GetFiles() {
// 		newLine := newFileTreeLine(file, depth+1)
// 		fileTree = append(fileTree, newLine)
// 	}
//
// 	return fileTree
// }

//	func styleFileTreeLine(file FileTreeItem) lipgloss.Style {
//		style := lipgloss.NewStyle()
//		if file.IsFullyStaged() {
//			style = style.Foreground(lipgloss.Color("#a6e3a1"))
//		} else {
//			style = style.Foreground(lipgloss.Color("#f38ba8"))
//		}
//
//		return style
//	}
//
//	func (ft FileTree) Update(msg tea.Msg) (FileTree, tea.Cmd) {
//		var cmd tea.Cmd
//		ft, cmd = ft.updateAsModel(msg)
//		return ft, cmd
//	}
//
//	func (ft FileTree) updateAsModel(msg tea.Msg) (FileTree, tea.Cmd) {
//		var cmd tea.Cmd
//		if !ft.IsFocused {
//			return ft, cmd
//		}
//
//		downKeymap := key.NewBinding(
//			key.WithKeys("down", "j"),
//			key.WithHelp("down/j", "Down"),
//		)
//
//		upKeymap := key.NewBinding(
//			key.WithKeys("up", "k"),
//			key.WithHelp("up/k", "Up"),
//		)
//
//		enterKeymap := key.NewBinding(
//			key.WithKeys("enter"),
//			key.WithHelp("enter", "Enter"),
//		)
//
//		spaceKeymap := key.NewBinding(
//			key.WithKeys(" "),
//			key.WithHelp("space", "Space"),
//		)
//
//		cKeymap := key.NewBinding(
//			key.WithKeys("c"),
//			key.WithHelp("c", "Commit"),
//		)
//
//		switch msg := msg.(type) {
//		case tea.KeyMsg:
//			switch {
//			case key.Matches(msg, downKeymap):
//				ft.cursorDown()
//			case key.Matches(msg, upKeymap):
//				ft.cursorUp()
//			case key.Matches(msg, enterKeymap):
//				ft.handleEnter()
//			case key.Matches(msg, spaceKeymap):
//				ft.handleSpace()
//			case key.Matches(msg, cKeymap):
//				ft.handleC()
//			}
//		}
//
//		return ft, cmd
//	}
//
//	func (ft *FileTree) handleC() {
//		hasStaged := false
//		for _, line := range ft.fileTreeLines {
//			if strings.HasPrefix(line.Item.GetStatus(), "M") {
//				hasStaged = true
//				break
//			}
//		}
//
//		if !hasStaged {
//			return
//		}
//
//		input := textinput.New("Commit message:")
//		input.Placeholder = ""
//
//		// filepath, err := input.RunPrompt()
//		// if err != nil {
//		// 	return
//		// }
//
//		// git.Commit(filepath)
//		ft.reloadModel()
//	}
//
//	func (ft *FileTree) handleSpace() {
//		// selectedLine, err := ft.getSelectedLine()
//		// if err != nil {
//		// 	return
//		// }
//
//		// filepath := selectedLine.Item.GetFilePath()
//
//		// if strings.HasSuffix(selectedLine.Item.GetStatus(), "M") {
//		// 	git.Stage(filepath)
//		// } else {
//		// 	git.Unstage(filepath)
//		// }
//
//		ft.reloadModel()
//	}
//
//	func (ft *FileTree) reloadModel() {
//		// rawStatus := git.GetRawStatus()
//		// directory := git.GetStatus(rawStatus)
//		// ft.fileTreeLines = buildFileTreeLines(directory)
//	}
//
//	func (ft *FileTree) handleEnter() {
//		selectedLine, err := ft.getSelectedLine()
//		if err != nil {
//			return
//		}
//
//		switch lineItem := selectedLine.Item.(type) {
//		case *git.Directory:
//			lineItem.ToggleExpanded()
//		}
//	}
//
//	func (ft *FileTree) cursorDown() {
//		for i := ft.cursorIndex + 1; i < len(ft.fileTreeLines); i++ {
//			if ft.updateCursorIndex(i) {
//				return
//			}
//		}
//	}
//
//	func (ft *FileTree) cursorUp() {
//		for i := ft.cursorIndex - 1; i >= 0; i-- {
//			if ft.updateCursorIndex(i) {
//				return
//			}
//		}
//	}
//
//	func (ft *FileTree) updateCursorIndex(newIndex int) bool {
//		newSelectedLine := ft.fileTreeLines[newIndex]
//		if newSelectedLine.Item.IsVisible() {
//			ft.cursorIndex = newIndex
//			return true
//		}
//
//		return false
//	}
//
//	func (ft FileTree) GetIndex() int {
//		return ft.cursorIndex
//	}
//
//	func (ft FileTree) getSelectedLine() (FileTreeLine, error) {
//		if len(ft.fileTreeLines) == 0 {
//			var result FileTreeLine
//			return result, errors.New("No file tree lines to display")
//		}
//
//		return ft.fileTreeLines[max(0, min(len(ft.fileTreeLines)-1, ft.cursorIndex))], nil
//	}
//
//	func (ft FileTree) GetSelectedFilepath() string {
//		currentLine, err := ft.getSelectedLine()
//		if err != nil {
//			return ""
//		}
//
//		return currentLine.Item.GetFilePath()
//	}
