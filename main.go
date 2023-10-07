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

	availableWidth := m.width
	columnWidth := availableWidth / 12
	diffWidth := columnWidth * 4

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		availableWidth := m.width
		columnWidth := availableWidth / 12
		diffWidth = columnWidth * 5

		height := msg.Height - styling.ColumnStyle.GetHeight() - 5

		if !m.ready {
			m.lviewport = viewport.New(diffWidth, height)

			m.rviewport = viewport.New(diffWidth, height)

			newFiletree, newCmd := m.fileTree.Update(msg)
			m.updateFiletree(newFiletree, diffWidth)
			cmds = append(cmds, newCmd)
			m.ready = true
			styling.DiffStyle.Width(diffWidth)
		} else {
			m.UpdateDiffDisplay(diffWidth)
			styling.DiffStyle.Width(diffWidth)

			m.lviewport.Width = diffWidth
			m.lviewport.Height = height

			m.rviewport.Width = diffWidth
			m.rviewport.Height = height
		}
	}

	newFiletree, newCmd := m.fileTree.Update(msg)
	m.updateFiletree(newFiletree, diffWidth)
	cmds = append(cmds, newCmd)

	m.lviewport, cmd = m.lviewport.Update(msg)
	cmds = append(cmds, cmd)

	m.rviewport, cmd = m.rviewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	takenWidth := m.lviewport.Width * 2
	headerStlying := styling.HeaderStyle.Width(m.width - 2)
	header := headerStlying.Render("Git-UI")

	leftDiff := styling.ColumnStyle.Render(m.lviewport.View())
	rightDiff := styling.ColumnStyle.Render(m.rviewport.View())

	width := m.width - takenWidth - styling.ColumnStyle.GetHorizontalBorderSize()
	fileTreeStyle := lipgloss.NewStyle().MaxWidth(width)
	fileTreeString := fileTreeStyle.Render(m.fileTree.Render())

	fileTree := styling.ColumnStyle.Copy().Width(width).Render(fileTreeString)

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
