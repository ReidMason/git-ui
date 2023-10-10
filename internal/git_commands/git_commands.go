package gitcommands

import (
	"fmt"
	"git-ui/internal/utils"
	"log"
)

type GitCommandRunner interface {
	Stage(filepath string)
	Unstage(filepath string)
	Commit(commitMessage string)
	GetDiff(filepath string) string
	GetStatus() string
}

type GitCommandLine struct{}

func New() GitCommandLine {
	return GitCommandLine{}
}

func (g GitCommandLine) Stage(filepath string) {
	_, err := utils.RunCommand("git", "add", "--", fmt.Sprintf(`%s`, filepath))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) Unstage(filepath string) {
	_, err := utils.RunCommand("git", "reset", "HEAD", "--", fmt.Sprintf(`%s`, filepath))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) Commit(commitMessage string) {
	_, err := utils.RunCommand("git", "commit", "-m", fmt.Sprintf(`%s`, commitMessage))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) GetDiff(filepath string) string {
	args := []string{"diff", "--no-ext-diff", "-U1000", "--"}
	// If it's new we want to add /dev/null instead

	args = append(args, filepath)

	result, err := utils.RunCommand("git", args...)

	if err != nil {
		log.Fatal("Failed to get git diff", err)
	}

	return result
}

func (g GitCommandLine) GetStatus() string {
	result, err := utils.RunCommand("git", "status", "-u", "--porcelain=v2", "--branch")

	if err != nil {
		log.Fatal("Failed to get git status")
	}

	return result
}
