package state

import (
	"git-ui/internal/git"
)

type State struct {
	gitStatus git.GitStatus
	message   string
	viewWidth uint16
}

func New() State {
	return State{}
}

func (s *State) SetGitStatus(gitStatus git.GitStatus) {
	s.gitStatus = gitStatus
}

func (s State) GetGitStatus() git.GitStatus { return s.gitStatus }

func (s *State) SetMessage(message string) {
	s.message = message
}

func (s State) GetMessage() string { return s.message }
