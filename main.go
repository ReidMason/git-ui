package main

import (
	"fmt"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/git"
	gitcommands "git-ui/internal/git_commands"
	"git-ui/internal/state"
	"git-ui/internal/styling"
	"git-ui/internal/ui"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App stores the main app information
// This includes all components like the file tree and the main app state
// The components get a pointer to the main app state so they can change it themselves
// So each component has it's own events and acts as a view but mutates the central state
//
// State is ONlY updated through setters
// Components or other "listeners" can subscribe to updates on the state which they can then act on
// This is how we will re-draw efficiently

type Model struct {
	state     state.State
	git       git.Git
	lviewport viewport.Model
	rviewport viewport.Model
	fileTree  filetree.FileTree
	ready     bool
}

func initModel() Model {
	gitCommands := gitcommands.New()

	model := Model{
		git:   git.New(gitCommands),
		state: state.New(),
	}

	gitStatus := model.git.GetStatus()
	model.state.SetGitStatus(gitStatus)
	model.fileTree = filetree.New(gitStatus.Directory)

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

// func (m Model) UpdateDiffDisplay(lineWidth int) {
// 	ldiff := styling.StyleDiff(m.diff.Diff1, lineWidth)
// 	m.lviewport.SetContent(ldiff)
//
// 	rdiff := styling.StyleDiff(m.diff.Diff2, lineWidth)
// 	m.rviewport.SetContent(rdiff)
// }

// func (m *Model) updateFiletree(newFiletree filetree.FileTree, lineWidth int) {
// 	if !m.ready || m.fileTree.GetIndex() != newFiletree.GetIndex() {
// 		filepath := newFiletree.GetSelectedFilepath()
// 		if filepath != "" {
// 			diffString := git.GetRawDiff(filepath)
// 			diff := git.ParseDiff(diffString)
// 			m.lviewport.SetContent(styling.StyleDiff(diff.Diff1, lineWidth))
// 			m.rviewport.SetContent(styling.StyleDiff(diff.Diff2, lineWidth))
// 			m.lviewport.GotoTop()
// 			m.rviewport.GotoTop()
// 			m.UpdateDiffDisplay(lineWidth)
// 		}
// 	}
//
// 	m.fileTree = newFiletree
// 	m.UpdateDiffDisplay(lineWidth)
// }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// availableWidth := m.width
	// columnWidth := availableWidth / 12
	// diffWidth := columnWidth * 4
	//
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 	rightKey := key.NewBinding(
		// 		key.WithKeys("right", "l"),
		// 		key.WithHelp("right/l", "Right"),
		// 	)
		//
		// 	leftKey := key.NewBinding(
		// 		key.WithKeys("left", "h"),
		// 		key.WithHelp("left/h", "Left"),
		// 	)
		//
		// 	switch {
		// 	case key.Matches(msg, rightKey):
		// 		m.isFocused = true
		// 		m.fileTree.IsFocused = false
		// 	case key.Matches(msg, leftKey):
		// 		m.isFocused = false
		// 		m.fileTree.IsFocused = true
		// 	}
		//
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			// m.state.SetMessage(msg.String())
		}
	case tea.WindowSizeMsg:
		m.state.SetViewWidth(msg.Width)
		// 	m.width = msg.Width
		// 	availableWidth := m.width
		// 	columnWidth := availableWidth / 12
		// 	diffWidth = columnWidth * 5
		//
		// 	height := msg.Height - styling.ColumnStyle.GetHeight() - 5
		//
		// 	if !m.ready {
		// 		m.lviewport = viewport.New(diffWidth, height)
		// 		m.rviewport = viewport.New(diffWidth, height)
		//
		// 		newFiletree, newCmd := m.fileTree.Update(msg)
		// 		m.updateFiletree(newFiletree, diffWidth)
		// 		cmds = append(cmds, newCmd)
		// 		m.ready = true
		// 		styling.DiffStyle.Width(diffWidth)
		// 	} else {
		// 		m.UpdateDiffDisplay(diffWidth)
		// 		styling.DiffStyle.Width(diffWidth)
		//
		// 		m.lviewport.Width = diffWidth
		// 		m.lviewport.Height = height
		//
		// 		m.rviewport.Width = diffWidth
		// 		m.rviewport.Height = height
		// 	}
		// }

		// if m.fileTree.IsFocused {
		// 	newFiletree, newCmd := m.fileTree.Update(msg)
		// 	m.updateFiletree(newFiletree, diffWidth)
		// 	cmds = append(cmds, newCmd)
	}

	// if m.isFocused {
	// 	m.lviewport, cmd = m.lviewport.Update(msg)
	// 	cmds = append(cmds, cmd)
	//
	// 	m.rviewport, cmd = m.rviewport.Update(msg)
	// 	cmds = append(cmds, cmd)
	// }

	m.fileTree.Update(msg)

	return m, tea.Batch(cmds...)
}

func formatDiff(viewport viewport.Model, focused bool) string {
	diffString := viewport.View()
	thing := lipgloss.
		NewStyle().
		MaxWidth(viewport.Width - styling.ColumnStyle.GetHorizontalBorderSize()).
		Render(diffString)

	borderColour := "60"
	if focused {
		borderColour = "62"
	}

	return styling.ColumnStyle.
		Copy().
		BorderForeground(lipgloss.Color(borderColour)).
		Width(viewport.Width - styling.ColumnStyle.GetHorizontalBorderSize()).
		Render(thing)
}

func (m Model) View() string {
	// takenWidth := m.lviewport.Width * 2
	//
	// leftDiff := formatDiff(m.lviewport, m.isFocused)
	// rightDiff := formatDiff(m.rviewport, m.isFocused)
	//
	// width := m.width - takenWidth - styling.ColumnStyle.GetHorizontalBorderSize()
	// fileTreeStyle := lipgloss.NewStyle().MaxWidth(width)
	// fileTreeString := fileTreeStyle.Render(m.fileTree.Render())
	//
	// borderColour := "60"
	// if m.fileTree.IsFocused {
	// 	borderColour = "62"
	// }
	//
	// fileTree := styling.ColumnStyle.
	// 	Copy().
	// 	BorderForeground(lipgloss.Color(borderColour)).
	// 	Width(width).
	// 	Render(fileTreeString)
	//
	// mainBody := lipgloss.JoinHorizontal(lipgloss.Left, fileTree, leftDiff, rightDiff)
	//

	return ui.RenderMainView(m.state.GetViewWidth(), m.fileTree)
}

func main() {
	debug := true
	if debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

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
