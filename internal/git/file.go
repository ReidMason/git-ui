package git

import (
	"fmt"
	"path/filepath"
	"strings"

	filetree "github.com/ReidMason/git-ui/internal/fileTree"
	"github.com/ReidMason/git-ui/internal/styling"

	"github.com/charmbracelet/lipgloss"
)

var (
	StagedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1"))

	UnstagedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8"))
)

type File struct {
	name           string
	secondName     string
	Parent         *Directory
	Dirpath        string
	IndexStatus    rune
	WorktreeStatus rune
}

func (f File) GetDisplay() string {
	indexStatus := string(f.IndexStatus)
	if f.IndexStatus != '.' {
		indexStatus = StagedStyle.Render(indexStatus)
	} else {
		indexStatus = UnstagedStyle.Render(indexStatus)
	}

	worktreeStatusAndName := fmt.Sprintf("%s %s", string(f.WorktreeStatus), f.getDisplayName())
	if f.WorktreeStatus == '.' {
		worktreeStatusAndName = StagedStyle.Render(worktreeStatusAndName)
	} else {
		worktreeStatusAndName = UnstagedStyle.Render(worktreeStatusAndName)
	}

	indexStatus = styling.TrimColourResetChar(indexStatus)

	return lipgloss.JoinHorizontal(0, indexStatus, worktreeStatusAndName)
}

func (f File) getDisplayName() string {
	if f.secondName != "" {
		return fmt.Sprintf("%s -> %s", f.secondName, f.name)
	}

	return f.name
}

func (f File) GetFilePath() string { return filepath.Join(f.Dirpath, f.name) }
func (f File) GetFilePaths() []string {
	filepaths := []string{filepath.Join(f.Dirpath, f.name)}
	if f.secondName != "" {
		filepaths = append(filepaths, filepath.Join(f.Dirpath, f.secondName))
	}
	return filepaths
}
func (f File) GetDirectories() []filetree.FileTreeItem { return []filetree.FileTreeItem{} }
func (f File) GetFiles() []filetree.FileTreeItem       { return []filetree.FileTreeItem{} }
func (f File) IsDirectory() bool                       { return false }
func (f File) IsStaged() bool                          { return f.WorktreeStatus == '.' }

func newFile(filePath string, indexStatus, workTreeStatus rune, secondName string) File {
	dirpath, filename := filepath.Split(filePath)
	dirpath = filepath.Clean(dirpath)

	return File{
		name:           filename,
		Dirpath:        dirpath,
		IndexStatus:    indexStatus,
		WorktreeStatus: workTreeStatus,
		secondName:     secondName,
	}
}

func addFile(directory *Directory, dirpath []string, visitedDirs []string, newFile File) {
	if len(dirpath) == 0 || dirpath[0] == "." {
		newFile.Parent = directory
		directory.Files = append(directory.Files, newFile)
		return
	}

	visitedDirs = append(visitedDirs, dirpath[0])

	for _, subdir := range directory.Directories {
		if subdir.Name == dirpath[0] {
			addFile(subdir, dirpath[1:], visitedDirs, newFile)
			return
		}
	}

	fullDirpath := strings.Join(visitedDirs, "/")
	newDir := newDirectory(dirpath[0], fullDirpath, directory)
	addFile(newDir, dirpath[1:], visitedDirs, newFile)
	newDir.Parent = directory
	directory.Directories = append(directory.Directories, newDir)
}
