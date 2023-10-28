package gitcommands

import (
	"fmt"
	"log"
	"strings"

	"github.com/ReidMason/git-ui/internal/utils"
)

type GitCommandRunner interface {
	Stage(filepath string)
	Unstage(filepath string)
	Commit(commitMessage string)
	GetDiff(filepath string) string
	GetStatus() (string, error)
}

const BaseCmd = "git"

type GitCommandLine struct{}

func New() GitCommandLine {
	return GitCommandLine{}
}

func (g GitCommandLine) Stage(filepath string) {
	_, err := runGitCommand("add", "--", fmt.Sprintf(`%s`, filepath))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) Unstage(filepath string) {
	_, err := runGitCommand("reset", "HEAD", "--", fmt.Sprintf(`%s`, filepath))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) Commit(commitMessage string) {
	_, err := utils.RunCommand(BaseCmd, "commit", "-m", fmt.Sprintf(`%s`, commitMessage))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) GetDiff(filepath string) string {
	result, err := runGitCommand("diff", "--no-ext-diff", "-U1000", "--", filepath)

	// If we got a blank result the file might be staged so we use cached
	if err == nil && result == "" {
		result, err = runGitCommand("diff", "--no-ext-diff", "--cached", "-U1000", "--", filepath)
	}

	// If we got a blank result the file probably isn't indexed so compare it to /dev/null
	if err == nil && result == "" {
		result, err = runGitCommand("diff", "--no-ext-diff", "-U1000", "--", "/dev/null", filepath)
	}

	if err != nil {
		return result
	}

	return result
}

func getRootDir() string {
	rootDir, err := utils.RunCommand(BaseCmd, "rev-parse", "--show-toplevel")
	if err != nil {
		return "."
	}

	return strings.TrimSpace(rootDir)
}

func runGitCommand(args ...string) (string, error) {
	rootDirectory := getRootDir()
	args = append([]string{"-C", rootDirectory}, args...)
	return utils.RunCommand(BaseCmd, args...)
}

func (g GitCommandLine) GetStatus() (string, error) {
	return runGitCommand("status", "-u", "--porcelain=v2", "--branch")
}
