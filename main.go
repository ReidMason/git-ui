package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ReidMason/git-ui/internal/git"
	gitcommands "github.com/ReidMason/git-ui/internal/git_commands"
	"github.com/ReidMason/git-ui/internal/state"
	"github.com/ReidMason/git-ui/internal/ui"

	filetree "github.com/ReidMason/git-ui/internal/fileTree"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type Model struct {
	textInput textinput.Model
	git       git.Git
	lviewport viewport.Model
	rviewport viewport.Model
	state     state.State
	fileTree  filetree.FileTree
	ready     bool
}

func initModel() Model {
	gitCommands := gitcommands.New()

	ti := textinput.New()
	ti.Placeholder = "Commit message"
	ti.CharLimit = 156

	model := Model{
		git:       git.New(gitCommands),
		state:     state.New(),
		textInput: ti,
	}

	model.fileTree.Focus()

	return model
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	gitStatus := m.git.GetStatus()
	updateStatusCmd := func() tea.Msg { return GitStatusUpdate{newGitStatus: gitStatus, oldFilepath: ""} }
	cmds = append(cmds, updateStatusCmd)

	firstFilepath := ""
	if gitStatus.Directory != nil {
		firstFilepath = gitStatus.Directory.GetFirstFilePath()
	}

	cmds = append(cmds, m.handleFileTreeChange(firstFilepath))

	return tea.Batch(cmds...)
}

type GitStatusUpdate struct {
	oldFilepath  string
	newGitStatus git.GitStatus
}

func (m Model) toggleStageFile() tea.Cmd {
	return func() tea.Msg {
		selectedElement, err := m.fileTree.GetSelectedItem()
		if err != nil {
			return nil
		}

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

type DiffUpdate struct {
	newDiff git.Diff
}

func (m *Model) handleFileTreeChange(filepath string) tea.Cmd {
	return func() tea.Msg {
		if filepath == "" {
			filepath = m.fileTree.GetSelectedFilepath()
		}
		diff := m.git.GetDiff(filepath)
		return DiffUpdate{newDiff: diff}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if !m.state.DiffsFocused() {
		m.fileTree, cmd = m.fileTree.Update(msg, m.toggleStageFile(), m.handleFileTreeChange(""))
		cmds = append(cmds, cmd)
	}

	if m.textInput.Focused() {
		m.textInput, _ = m.textInput.Update(msg)
	}

	if m.state.DiffsFocused() {
		m.lviewport, cmd = m.lviewport.Update(msg)
		cmds = append(cmds, cmd)
		m.rviewport, cmd = m.rviewport.Update(msg)
		cmds = append(cmds, cmd)
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
	case DiffUpdate:
		m.lviewport.SetContent(ui.DiffToString(msg.newDiff.Diff1))
		m.rviewport.SetContent(ui.DiffToString(msg.newDiff.Diff2))

		m.lviewport.GotoTop()
		m.rviewport.GotoTop()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.state = m.state.SetViewWidth(msg.Width)
	m.state.SetViewHeight(msg.Height)

	if !m.ready {
		diffWidth, diffHeight := ui.GetDiffDimensions(msg.Width, msg.Height)
		m.lviewport = viewport.New(diffWidth, diffHeight)
		m.rviewport = viewport.New(diffWidth, diffHeight)

		m.ready = true
	} else {
		diffWidth, diffHeight := ui.GetDiffDimensions(msg.Width, msg.Height)
		m.lviewport.Width = diffWidth
		m.lviewport.Height = diffHeight

		m.rviewport.Width = diffWidth
		m.rviewport.Height = diffHeight
	}

	if m.textInput.Focused() {
		m.textInput.Width = m.getTextInputWidth()
	}

	return m, nil
}

func (m Model) handleKeypress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "a":
		if m.fileTree.Focused() {
			return m, func() tea.Msg {
				if m.state.GetGitStatus().Directory.GetStagedStatus() == git.FullyStaged {
					m.git.Unstage(".")
				} else {
					m.git.Stage(".")
				}

				return GitStatusUpdate{
					newGitStatus: m.git.GetStatus(),
					oldFilepath:  m.state.SelectedFilepath(),
				}
			}
		}
	case "enter":
		if m.textInput.Focused() {
			m.textInput.Blur()
			m.fileTree.Focus()
			commitMessage := m.textInput.Value()
			m.textInput.SetValue("")

			return m, func() tea.Msg {
				m.git.Commit(commitMessage)
				return GitStatusUpdate{
					newGitStatus: m.git.GetStatus(),
					oldFilepath:  m.state.SelectedFilepath(),
				}
			}
		} else if m.fileTree.Focused() {
			selectedItem, err := m.fileTree.GetSelectedItem()
			if err != nil {
				return m, nil
			}
			if !selectedItem.GetItem().IsDirectory() {
				m.fileTree.Blur()
				m.state = m.state.SetDiffsFocused(true)
			}
		}
	case "esc":
		if m.state.DiffsFocused() {
			m.fileTree.Focus()
			m.state = m.state.SetDiffsFocused(false)
		} else {
			m.textInput.Blur()
			m.fileTree.Focus()

			return m, func() tea.Msg {
				return GitStatusUpdate{
					newGitStatus: m.git.GetStatus(),
					oldFilepath:  m.state.SelectedFilepath(),
				}
			}
		}
	case "c":
		if !m.textInput.Focused() {
			m.fileTree.Blur()
			m.textInput.Focus()
			m.state = m.state.SetDiffsFocused(false)
			m.textInput.Width = m.getTextInputWidth()
		}
	}

	return m, nil
}

func (m Model) getTextInputWidth() int {
	statusbarText := ui.GetFooterTextContent(m.state)
	outputLength := lipgloss.Width(statusbarText)
	rightPadding := 6
	return m.state.ViewWidth() - outputLength - rightPadding
}

func (m Model) View() string {
	leftDiff := m.lviewport.View()
	rightDiff := m.rviewport.View()

	divider := ""
	if len(strings.TrimSpace(leftDiff)) > 0 {
		divider = strings.Repeat("│\n", lipgloss.Height(leftDiff))
	}

	diffsStyle := ui.GetDiffsDividerStyle(m.state)
	divider = diffsStyle.Render(strings.TrimSuffix(divider, "\n"))
	diffs := lipgloss.JoinHorizontal(0, leftDiff, divider, rightDiff)

	statusBar := ui.RenderStatusBar(m.state, m.textInput)
	display := ui.RenderMainView(m.fileTree, diffs, statusBar, m.state)

	return display
}

func main() {
	debug := false
	err := godotenv.Load()
	if err == nil {
		debug = "true" == os.Getenv("DEBUG_GIT_UI")
	}

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
