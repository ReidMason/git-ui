package state

import (
	"git-ui/internal/git"
)

type State struct {
	selectedFilepath string
	gitStatus        git.GitStatus
	viewWidth        int
	viewHeight       int
}

func New() State {
	return State{}
}

func (s State) GetGitStatus() git.GitStatus { return s.gitStatus }
func (s State) SetGitStatus(gitStatus git.GitStatus) State {
	s.gitStatus = gitStatus
	return s
}

func (s State) GetViewWidth() int           { return s.viewWidth }
func (s *State) SetViewWidth(viewWidth int) { s.viewWidth = viewWidth }

func (s State) GetViewHeight() int            { return s.viewHeight }
func (s *State) SetViewHeight(viewHeight int) { s.viewHeight = viewHeight }

func (s State) GetSelectedFilepath() string                  { return s.selectedFilepath }
func (s *State) SetSelectedFilepath(selectedFilepath string) { s.selectedFilepath = selectedFilepath }
