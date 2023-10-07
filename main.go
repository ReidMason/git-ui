package main

import (
	"fmt"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/git"
	"git-ui/internal/styling"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	diff      git.Diff
	lviewport viewport.Model
	rviewport viewport.Model

	gitStatus git.Directory
	fileTree  filetree.FileTree

	width int
	ready bool
}

func initModel() Model {
	rawStatus := git.GetRawStatus()
	// 	rawStatus = `# branch.oid d2ca38080c8a408cd3b4824d64237b8875acf98e
	// # branch.head main
	// # branch.upstream origin/main
	// # branch.ab +0 -0
	// 1 .M N... 100644 100644 100644 0a46fa61504b531ca53005dee44cc0b1cd6ffc99 0a46fa61504b531ca53005dee44cc0b1cd6ffc99 internal/fileTree/fileTree.go
	// 1 .M N... 100644 100644 100644 2107e7a49e44e0f97915bb523729889d9578a612 2107e7a49e44e0f97915bb523729889d9578a612 internal/git/git.go
	// 1 .M N... 100644 100644 100644 eba8c5554a26db93531ce2c90b34da40f86f887f eba8c5554a26db93531ce2c90b34da40f86f887f internal/git/git_test.go
	// 1 .M N... 100644 100644 100644 c789db6decaa4c7af3d5eb2214aea59f430dd5b1 c789db6decaa4c7af3d5eb2214aea59f430dd5b1 internal/utils/utils.go
	// 1 .M N... 100644 100644 100644 587b38dd887ed1cdd4fd9f45819f1e9f9d3ceca6 587b38dd887ed1cdd4fd9f45819f1e9f9d3ceca6 main.go`
	gitStatus := git.GetStatus(rawStatus)
	fileTree := filetree.New(gitStatus)

	return Model{
		gitStatus: gitStatus,
		ready:     false,
		fileTree:  fileTree,
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) UpdateDiffDisplay(lineWidth int) {
	ldiff := styling.StyleDiff(m.diff.Diff1, lineWidth)
	m.lviewport.SetContent(ldiff)

	rdiff := styling.StyleDiff(m.diff.Diff2, lineWidth)
	m.rviewport.SetContent(rdiff)
}

func (m *Model) updateFiletree(newFiletree filetree.FileTree, lineWidth int) {
	if !m.ready || m.fileTree.GetIndex() != newFiletree.GetIndex() {
		filepath := newFiletree.GetSelectedFilepath()
		if filepath != "" {
			diffString := git.GetRawDiff(filepath)
			diff := git.GetDiff(diffString)
			m.lviewport.SetContent(styling.StyleDiff(diff.Diff1, lineWidth))
			m.rviewport.SetContent(styling.StyleDiff(diff.Diff2, lineWidth))
			m.UpdateDiffDisplay(lineWidth)
		}
	}

	m.fileTree = newFiletree
	m.UpdateDiffDisplay(lineWidth)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	lineWidth := m.width - styling.ColumnStyle.GetHorizontalPadding()

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
		lineWidth = width - styling.ColumnStyle.GetHorizontalPadding()
		height := msg.Height - styling.ColumnStyle.GetVerticalPadding() - 5

		if !m.ready {

			m.lviewport = viewport.New(width, height)
			m.lviewport.YPosition = 10

			m.rviewport = viewport.New(width, height)
			m.rviewport.YPosition = 10

			newFiletree, newCmd := m.fileTree.Update(msg)
			m.updateFiletree(newFiletree, lineWidth)
			cmds = append(cmds, newCmd)

			styling.ColumnStyle.Width(width)
			m.ready = true
		} else {
			styling.ColumnStyle.Width(width)
			styling.ColumnStyle.Height(height)

			m.UpdateDiffDisplay(lineWidth)

			m.lviewport.Width = width
			m.lviewport.Height = height

			m.rviewport.Width = width
			m.rviewport.Height = height
		}
	}

	newFiletree, newCmd := m.fileTree.Update(msg)
	m.updateFiletree(newFiletree, lineWidth)
	cmds = append(cmds, newCmd)

	m.lviewport, cmd = m.lviewport.Update(msg)
	cmds = append(cmds, cmd)

	m.rviewport, cmd = m.rviewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	headerStlying := styling.HeaderStyle.Width(m.width - 10)
	header := headerStlying.Render("Git diff")

	fileTree := m.fileTree.Render()
	leftDiff := styling.ColumnStyle.Render(m.lviewport.View())
	rightDiff := styling.ColumnStyle.Render(m.rviewport.View())

	mainBody := lipgloss.JoinHorizontal(lipgloss.Left, fileTree, leftDiff, rightDiff)

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
