package state

import (
	"github.com/ReidMason/git-ui/internal/git"
)

type State struct {
	selectedFilepath string
	gitStatus        git.GitStatus
	viewWidth        int
	viewHeight       int
	diffsFocused     bool
}

func New() State {
	return State{}
}

func (s State) GetGitStatus() git.GitStatus { return s.gitStatus }
func (s State) SetGitStatus(gitStatus git.GitStatus) State {
	s.gitStatus = gitStatus
	return s
}

func (s State) SetDiffsFocused(focused bool) State {
	s.diffsFocused = focused
	return s
}
func (s State) DiffsFocused() bool { return s.diffsFocused }

func (s State) ViewWidth() int { return s.viewWidth }
func (s State) SetViewWidth(viewWidth int) State {
	s.viewWidth = viewWidth
	return s
}

func (s State) ViewHeight() int               { return s.viewHeight }
func (s *State) SetViewHeight(viewHeight int) { s.viewHeight = viewHeight }

func (s State) SelectedFilepath() string                     { return s.selectedFilepath }
func (s *State) SetSelectedFilepath(selectedFilepath string) { s.selectedFilepath = selectedFilepath }
