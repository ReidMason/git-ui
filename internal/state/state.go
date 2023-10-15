package state

import (
	"git-ui/internal/git"
)

type State struct {
	gitStatus  git.GitStatus
	message    string
	viewWidth  int
	viewHeight int
}

func New() State {
	return State{}
}

func (s State) GetGitStatus() git.GitStatus           { return s.gitStatus }
func (s *State) SetGitStatus(gitStatus git.GitStatus) { s.gitStatus = gitStatus }

func (s State) GetMessage() string         { return s.message }
func (s *State) SetMessage(message string) { s.message = message }

func (s State) GetViewWidth() int           { return s.viewWidth }
func (s *State) SetViewWidth(viewWidth int) { s.viewWidth = viewWidth }

func (s State) GetViewHeight() int            { return s.viewHeight }
func (s *State) SetViewHeight(viewHeight int) { s.viewHeight = viewHeight }
