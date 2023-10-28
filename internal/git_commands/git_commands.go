package gitcommands

import (
	"fmt"
	"git-ui/internal/utils"
	"log"
	"strings"
)

type GitCommandRunner interface {
	Stage(filepath string)
	Unstage(filepath string)
	Commit(commitMessage string)
	GetDiff(filepath string) string
	GetStatus() (string, error)
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
	// If it's new we want to add /dev/null instead
	result, err := runGitCommand("diff", "--no-ext-diff", "-U1000", "--", filepath)

	if err != nil {
		return ""
	}

	return result
}

func getRootDir() string {
	rootDir, err := utils.RunCommand("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "."
	}

	return strings.TrimSpace(rootDir)
}

func runGitCommand(args ...string) (string, error) {
	rootDirectory := getRootDir()
	args = append([]string{"-C", rootDirectory}, args...)
	return utils.RunCommand("git", args...)
}

func (g GitCommandLine) GetStatus() (string, error) {
	return runGitCommand("status", "-u", "--porcelain=v2", "--branch")
}
