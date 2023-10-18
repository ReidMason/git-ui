package main

import (
	"fmt"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/git"
	gitcommands "git-ui/internal/git_commands"
	"git-ui/internal/state"
	"git-ui/internal/ui"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
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

// BUG: If you have a file selected and the filetree refreshes and that item vanishes the selection gets broken
// Write a test for this

type Model struct {
	textInput        textinput.Model
	git              git.Git
	selectedFilepath string
	lviewport        viewport.Model
	rviewport        viewport.Model
	diff             git.Diff
	state            state.State
	fileTree         filetree.FileTree
	committing       bool
	ready            bool
}

func initModel() Model {
	gitCommands := gitcommands.New()

	model := Model{
		git:        git.New(gitCommands),
		state:      state.New(),
		committing: false,
	}

	gitStatus := model.git.GetStatus()
	model.state.SetGitStatus(gitStatus)
	model.fileTree = filetree.New(gitStatus.Directory)

	model.selectedFilepath = model.fileTree.GetSelectedFilepath()
	if model.selectedFilepath != "" {
		model.diff = model.git.GetDiff(model.selectedFilepath)
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	selectedItem := m.fileTree.Update(msg)
	if selectedItem != nil {
		filepath := selectedItem.GetFilePath()

		stage := true
		switch item := selectedItem.(type) {
		case *git.Directory:
			stage = item.GetStagedStatus() != git.FullyStaged
		case git.File:
			stage = !item.IsStaged()
		}

		if stage {
			m.git.Stage(filepath)
		} else {
			m.git.Unstage(filepath)
		}

		newStatus := m.git.GetStatus()
		m.state.SetGitStatus(newStatus)
		m.fileTree.UpdateDirectoryTree(newStatus.Directory, filepath)
	}

	if m.committing {
		m.textInput, _ = m.textInput.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.committing {
				m.committing = false
				m.fileTree.SetFocused(true)
				commitMessage := m.textInput.Value()
				m.git.Commit(commitMessage)

				newStatus := m.git.GetStatus()
				m.state.SetGitStatus(newStatus)
				m.fileTree.UpdateDirectoryTree(newStatus.Directory, m.selectedFilepath)
			}
		case "esc":
			m.committing = false
			m.fileTree.SetFocused(true)

			newStatus := m.git.GetStatus()
			m.state.SetGitStatus(newStatus)
			m.fileTree.UpdateDirectoryTree(newStatus.Directory, m.selectedFilepath)
		case "c":
			if !m.committing {
				m.committing = true
				m.fileTree.SetFocused(false)
				ti := textinput.New()
				ti.Placeholder = "Commit message"
				ti.Focus()
				ti.CharLimit = 156

				// This needs to be three less than it's actual width to account for the extra characters
				ti.Width = 47 // m.state.GetViewWidth() - 10

				m.textInput = ti
			}
		default:
			// m.state.SetMessage(msg.String())
		}

	case tea.WindowSizeMsg:
		m.state.SetViewWidth(msg.Width)
		m.state.SetViewHeight(msg.Height)

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

	newSelectedFilepath := m.fileTree.GetSelectedFilepath()
	if newSelectedFilepath != m.selectedFilepath {
		m.selectedFilepath = newSelectedFilepath
		m.diff = m.git.GetDiff(newSelectedFilepath)
		m.lviewport.SetContent(git.DiffToString(m.diff.Diff1))
		m.rviewport.SetContent(git.DiffToString(m.diff.Diff2))
		m.lviewport.GotoTop()
		m.rviewport.GotoTop()
	}

	// m.lviewport, cmd = m.lviewport.Update(msg)
	// cmds = append(cmds, cmd)
	// m.rviewport, cmd = m.rviewport.Update(msg)
	// cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	width := m.state.GetViewWidth()
	height := m.state.GetViewHeight()
	leftDiff := m.lviewport.View()
	rightDiff := m.rviewport.View()
	diffs := lipgloss.JoinHorizontal(0, leftDiff, rightDiff)

	statusBar := ui.RenderStatusBar(m.state.GetGitStatus(), width, m.textInput.View(), m.committing)
	display := ui.RenderMainView(width, height, m.fileTree, diffs, statusBar)

	return display
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
