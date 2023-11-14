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

type GitCommandLine struct {
	rootDir string
}

func New(rootDir string) GitCommandLine {
	return GitCommandLine{
		rootDir: rootDir,
	}
}

func (g GitCommandLine) Stage(filepath string) {
	_, err := g.runGitCommand("add", "--", fmt.Sprintf(`%s`, filepath))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) Unstage(filepath string) {
	_, err := g.runGitCommand("reset", "HEAD", "--", fmt.Sprintf(`%s`, filepath))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) Commit(commitMessage string) {
	_, err := g.runGitCommand("commit", "-m", fmt.Sprintf(`%s`, commitMessage))

	if err != nil {
		log.Printf("Got err: %s", err)
	}
}

func (g GitCommandLine) GetDiff(filepath string) string {
	result, err := g.runGitCommand("diff", "--no-ext-diff", "-U1000", "--", filepath)

	// If we got a blank result the file might be staged so we use cached
	if err == nil && result == "" {
		result, err = g.runGitCommand("diff", "--no-ext-diff", "--cached", "-U1000", "--", filepath)
	}

	// If we got a blank result the file probably isn't indexed so compare it to /dev/null
	if err == nil && result == "" {
		result, err = g.runGitCommand("diff", "--no-ext-diff", "-U1000", "--", "/dev/null", filepath)
	}

	if err != nil {
		return result
	}

	return result
}

func GetRootDir() string {
	rootDir, err := utils.RunCommand(BaseCmd, "rev-parse", "--show-toplevel")
	if err != nil {
		return "."
	}

	return strings.TrimSpace(rootDir)
}

func (g GitCommandLine) runGitCommand(args ...string) (string, error) {
	args = append([]string{"-C", g.rootDir}, args...)
	return utils.RunCommand(BaseCmd, args...)
}

func (g GitCommandLine) GetStatus() (string, error) {
	return g.runGitCommand("status", "-u", "--porcelain=v2", "--branch")
}
