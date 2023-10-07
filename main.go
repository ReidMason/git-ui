package main

import (
	"fmt"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/git"
	"git-ui/internal/styling"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ldiff []git.DiffLine

	rdiff     []git.DiffLine
	lviewport viewport.Model
	rviewport viewport.Model

	gitStatus git.Directory
	fileLine  int
	fileTree  []filetree.FileTreeLine

	width int
	ready bool
}

func initModel() Model {
	rawStatus := git.GetRawStatus()
	rawStatus = `# branch.oid c86e7ed35f16570194c2308a2f8cb53155d0440d
# branch.head main
# branch.upstream origin/main
# branch.ab +0 -0
1 .M N... 100644 100644 100644 51d742a142700c40e5d5d4915b44da5d238bef81 51d742a142700c40e5d5d4915b44da5d238bef81 internal/git/git.go
1 .M N... 100644 100644 100644 8508f049bcb61d4c52d92e5a4c9a71051f00bcba 8508f049bcb61d4c52d92e5a4c9a71051f00bcba internal/git/git_test.go
1 M. N... 100644 100644 100644 1cdd739f6591c3aca07eab977748142a1ba14056 c345bc6f17650da4f51350e8faa56e4f4c61663e main.go
? internal/styling/styling.go`
	gitStatus := git.GetStatus(rawStatus)
	fileTree := buildFileTree(&gitStatus, make([]filetree.FileTreeLine, 0), -1)

	return Model{
		gitStatus: gitStatus,
		ready:     false,
		fileLine:  0,
		fileTree:  fileTree,
	}
}

func buildFileTree(directory *git.Directory, fileTree []filetree.FileTreeLine, depth int) []filetree.FileTreeLine {
	newLine := filetree.New(directory, depth)
	fileTree = append(fileTree, newLine)

	for _, subDirectory := range directory.Directories {
		fileTree = buildFileTree(&subDirectory, fileTree, depth+1)
	}

	for _, file := range directory.Files {
		newLine := filetree.New(file, depth+1)
		fileTree = append(fileTree, newLine)
	}

	return fileTree
}

func buildFileTreeString(fileTree []filetree.FileTreeLine) []string {
	output := make([]string, 0)
	for i := 1; i < len(fileTree); i++ {
		line := fileTree[i]

		prefix := strings.Repeat(" ", line.Depth) + "-"
		lineString := prefix + line.Item.GetStatus() + " " + line.Item.GetName()
		style := styling.StyleFileTreeLine(line.Item)
		output = append(output, style.Render(lineString))

		if !line.Item.IsExpanded() {
			i += line.Item.Children()
		}
	}

	return output
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		offset := 2
		m.width = msg.Width
		width := m.width/2 - offset
		lineWidth := width - styling.ColumnStyle.GetHorizontalPadding()
		height := msg.Height - styling.ColumnStyle.GetVerticalPadding() - 5

		if !m.ready {
			lines := buildFileTreeString(m.fileTree)
			fs := filetree.Render(lines)

			m.lviewport = viewport.New(width, height)
			m.lviewport.YPosition = 10
			// ldiff := styleDiff(m.ldiff, lineWidth)
			m.lviewport.SetContent(fs)

			m.rviewport = viewport.New(width, height)
			m.rviewport.YPosition = 10

			rdiff := styling.StyleDiff(m.rdiff, lineWidth)
			m.rviewport.SetContent(rdiff)

			styling.ColumnStyle.Width(width)
			m.ready = true
		} else {
			styling.ColumnStyle.Width(width)
			styling.ColumnStyle.Height(height)

			// fs := ""
			// for _, file := range m.gitStatus {
			// 	fs += file.Name + "\n"
			// }

			// ldiff := styleDiff(m.ldiff, lineWidth)
			// m.lviewport.SetContent(fs)

			rdiff := styling.StyleDiff(m.rdiff, lineWidth)
			m.rviewport.SetContent(rdiff)

			m.lviewport.Width = width
			m.lviewport.Height = height

			m.rviewport.Width = width
			m.rviewport.Height = height
		}
	}

	m.lviewport, cmd = m.lviewport.Update(msg)
	cmds = append(cmds, cmd)

	m.rviewport, cmd = m.rviewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	headerStlying := styling.HeaderStyle.Width(m.width - 2)
	header := headerStlying.Render("Git diff")

	leftView := styling.ColumnStyle.Render(m.lviewport.View())
	rightView := styling.ColumnStyle.Render(m.rviewport.View())

	mainBody := lipgloss.JoinHorizontal(lipgloss.Left, leftView, rightView)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainBody)
}

func main() {
	p := tea.NewProgram(
		initModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
