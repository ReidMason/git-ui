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
	args := []string{"diff", "--no-ext-diff", "-U1000", "--"}
	// If it's new we want to add /dev/null instead

	args = append(args, filepath)

	result, err := utils.RunCommand("git", args...)

	if err != nil {
		return ""
	}

	return result
}

func getRootDir() (string, error) {
	rootDir, err := utils.RunCommand("git", "rev-parse", "--show-toplevel")
	return strings.TrimSpace(rootDir), err
}

func (g GitCommandLine) GetStatus() (string, error) {
	rootDirectory, err := getRootDir()
	if err != nil {
		log.Println("Failed to get rootDirectory")
		rootDirectory = "."
	}

	log.Println("git", "-C", rootDirectory, "status", "-u", "--porcelain=v2", "--branch")

	return utils.RunCommand("git", "-C", rootDirectory, "status", "-u", "--porcelain=v2", "--branch")
}
