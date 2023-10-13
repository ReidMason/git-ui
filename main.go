package main

import (
	"fmt"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/git"
	gitcommands "git-ui/internal/git_commands"
	"git-ui/internal/state"
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
	state            state.State
	git              git.Git
	lviewport        viewport.Model
	rviewport        viewport.Model
	selectedFilepath string
	diff             git.Diff
	fileTree         filetree.FileTree
	ready            bool
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

	model.selectedFilepath = model.fileTree.GetSelectedFilepath()
	model.diff = model.git.GetDiff(model.selectedFilepath)

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
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			// m.state.SetMessage(msg.String())
		}

	case tea.WindowSizeMsg:
		m.state.SetViewWidth(msg.Width)

		if !m.ready {
			diffWidth, diffHeight := ui.GetDiffDimensions(msg.Width, msg.Height)
			m.lviewport = viewport.New(diffWidth, diffHeight)
			m.rviewport = viewport.New(diffWidth, diffHeight)
			m.lviewport.SetContent(git.DiffToString(m.diff.Diff1))
			m.rviewport.SetContent(git.DiffToString(m.diff.Diff2))

			m.ready = true
		} else {
			diffWidth, diffHeight := ui.GetDiffDimensions(msg.Width, msg.Height)
			m.lviewport.Width = diffWidth
			m.lviewport.Height = diffHeight

			m.rviewport.Width = diffWidth
			m.rviewport.Height = diffHeight
		}
	}

	m.lviewport, cmd = m.lviewport.Update(msg)
	cmds = append(cmds, cmd)
	m.rviewport, cmd = m.rviewport.Update(msg)
	cmds = append(cmds, cmd)

	m.fileTree.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	width := m.state.GetViewWidth()
	leftDiff := m.lviewport.View()
	rightDiff := m.rviewport.View()
	diffs := lipgloss.JoinHorizontal(0, leftDiff, rightDiff)
	return ui.RenderMainView(width, m.fileTree, diffs)
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
