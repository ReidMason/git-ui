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
	textInput  textinput.Model
	git        git.Git
	lviewport  viewport.Model
	rviewport  viewport.Model
	state      state.State
	fileTree   filetree.FileTree
	committing bool
	ready      bool
}

func initModel() Model {
	gitCommands := gitcommands.New()

	model := Model{
		git:        git.New(gitCommands),
		state:      state.New(),
		committing: false,
	}

	model.fileTree.SetFocused(true)

	return model
}

func (m Model) Init() tea.Cmd {
	gitStatus := m.git.GetStatus()
	return func() tea.Msg { return GitStatusUpdate{newGitStatus: gitStatus, oldFilepath: ""} }
}

type GitStatusUpdate struct {
	oldFilepath  string
	newGitStatus git.GitStatus
}

func toggleStageFile(m Model) tea.Cmd {
	return func() tea.Msg {
		selectedElement := m.fileTree.GetSelectedItem()
		selectedItem := selectedElement.GetItem()
		if selectedItem == nil {
			return nil
		}

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

		return GitStatusUpdate{newGitStatus: m.git.GetStatus(), oldFilepath: filepath}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.fileTree, cmd = m.fileTree.Update(msg, toggleStageFile(m))
	cmds = append(cmds, cmd)

	if m.committing {
		m.textInput, _ = m.textInput.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd = m.handleKeypress(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m, cmd = m.handleWindowSizeMsg(msg)
		cmds = append(cmds, cmd)
	case GitStatusUpdate:
		m.state = m.state.SetGitStatus(msg.newGitStatus)
		m.fileTree = m.fileTree.UpdateDirectoryTree(msg.newGitStatus.Directory, msg.oldFilepath)
	}

	newSelectedFilepath := m.fileTree.GetSelectedFilepath()
	if newSelectedFilepath != m.state.GetSelectedFilepath() {
		m.state.SetSelectedFilepath(newSelectedFilepath)
		m.state.SetDiff(m.git.GetDiff(newSelectedFilepath))
		m.lviewport.SetContent(git.DiffToString(m.state.GetDiff().Diff1))
		m.rviewport.SetContent(git.DiffToString(m.state.GetDiff().Diff2))
		m.lviewport.GotoTop()
		m.rviewport.GotoTop()
	}

	// m.lviewport, cmd = m.lviewport.Update(msg)
	// cmds = append(cmds, cmd)
	// m.rviewport, cmd = m.rviewport.Update(msg)
	// cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.state.SetViewWidth(msg.Width)
	m.state.SetViewHeight(msg.Height)

	if !m.ready {
		diffWidth, diffHeight := ui.GetDiffDimensions(msg.Width, msg.Height)
		m.lviewport = viewport.New(diffWidth, diffHeight)
		m.rviewport = viewport.New(diffWidth, diffHeight)
		m.lviewport.SetContent(git.DiffToString(m.state.GetDiff().Diff1))
		m.rviewport.SetContent(git.DiffToString(m.state.GetDiff().Diff2))

		m.ready = true
	} else {
		diffWidth, diffHeight := ui.GetDiffDimensions(msg.Width, msg.Height)
		m.lviewport.Width = diffWidth
		m.lviewport.Height = diffHeight

		m.rviewport.Width = diffWidth
		m.rviewport.Height = diffHeight
	}

	return m, nil
}

func (m Model) handleKeypress(msg tea.KeyMsg) (Model, tea.Cmd) {
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
			m.fileTree = m.fileTree.UpdateDirectoryTree(newStatus.Directory, m.state.GetSelectedFilepath())
		}
	case "esc":
		m.committing = false
		m.fileTree.SetFocused(true)

		newStatus := m.git.GetStatus()
		m.state.SetGitStatus(newStatus)
		m.fileTree = m.fileTree.UpdateDirectoryTree(newStatus.Directory, m.state.GetSelectedFilepath())
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
	}

	return m, nil
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
