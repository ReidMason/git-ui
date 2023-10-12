package git

import (
	"fmt"
	filetree "git-ui/internal/fileTree"
	"git-ui/internal/styling"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

var (
	StagedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1"))

	UnstagedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8"))
)

type File struct {
	Name           string
	Parent         *Directory
	Dirpath        string
	IndexStatus    rune
	WorktreeStatus rune
}

func (f File) GetDisplay() string {
	indexStatus := string(f.IndexStatus)
	if f.IndexStatus == 'M' {
		indexStatus = StagedStyle.Render(indexStatus)
	} else {
		indexStatus = UnstagedStyle.Render(indexStatus)
	}

	worktreeStatusAndName := fmt.Sprintf("%s %s", string(f.WorktreeStatus), f.Name)
	if f.WorktreeStatus == '.' {
		worktreeStatusAndName = StagedStyle.Render(worktreeStatusAndName)
	} else {
		worktreeStatusAndName = UnstagedStyle.Render(worktreeStatusAndName)
	}

	indexStatus = styling.TrimColourResetChar(indexStatus)

	return lipgloss.JoinHorizontal(0, indexStatus, worktreeStatusAndName)
}
func (f File) GetFilePath() string                     { return filepath.Join(f.Dirpath, f.Name) }
func (f File) GetDirectories() []filetree.FileTreeItem { return []filetree.FileTreeItem{} }
func (f File) GetFiles() []filetree.FileTreeItem       { return []filetree.FileTreeItem{} }
func (f File) IsDirectory() bool                       { return false }

func newFile(filePath string, indexStatus, workTreeStatus rune) File {
	dirpath, filename := filepath.Split(filePath)
	dirpath = filepath.Clean(dirpath)

	return File{Name: filename, Dirpath: dirpath, IndexStatus: indexStatus, WorktreeStatus: workTreeStatus}
}

func addFile(directory *Directory, dirpath []string, newFile File) {
	if len(dirpath) == 0 || dirpath[0] == "." {
		newFile.Parent = directory
		directory.Files = append(directory.Files, newFile)
		return
	}

	for _, subdir := range directory.Directories {
		if subdir.Name == dirpath[0] {
			addFile(subdir, dirpath[1:], newFile)
			// directory.Directories[i] = addFile(subdir, dirpath[1:], newFile)
			return
		}
	}

	newDir := newDirectory(dirpath[0], directory)
	addFile(newDir, dirpath[1:], newFile)
	newDir.Parent = directory
	directory.Directories = append(directory.Directories, newDir)
}
