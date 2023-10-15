package git

import (
	"git-ui/internal/colours"
	filetree "git-ui/internal/fileTree"

	"github.com/charmbracelet/lipgloss"
)

type StagedStatus int

const (
	Unstaged StagedStatus = iota
	PartiallyStaged
	FullyStaged
)

type Directory struct {
	Name        string
	Filepath    string
	Parent      *Directory
	Files       []File
	Directories []*Directory
}

func newDirectory(name, filepath string, parent *Directory) *Directory {
	return &Directory{
		Name:        name,
		Filepath:    filepath,
		Files:       make([]File, 0),
		Directories: make([]*Directory, 0),
		Parent:      parent,
	}
}

func (d Directory) GetDirectories() []filetree.FileTreeItem {
	items := make([]filetree.FileTreeItem, len(d.Directories))
	for i, dir := range d.Directories {
		items[i] = dir
	}
	return items
}

func (d Directory) GetFiles() []filetree.FileTreeItem {
	items := make([]filetree.FileTreeItem, len(d.Files))
	for i, file := range d.Files {
		items[i] = file
	}
	return items
}

func (d Directory) GetDisplay() string {
	status := d.getStagedStatus()
	styling := lipgloss.NewStyle()
	text := d.Name

	if status == FullyStaged {
		return styling.Foreground(lipgloss.Color(colours.Green)).Render(text)
	}

	if status == PartiallyStaged {
		return styling.Foreground(lipgloss.Color(colours.Peach)).Render(text)
	}

	return styling.Foreground(lipgloss.Color(colours.Red)).Render(text)
}

func (d Directory) getStagedStatus() StagedStatus {
	for _, subDirectory := range d.Directories {
		subDirectoryStatus := subDirectory.getStagedStatus()
		if subDirectoryStatus != FullyStaged {
			return subDirectoryStatus
		}
	}

	if len(d.Files) == 0 {
		return FullyStaged
	}

	hasStagedFile := false
	hasUnstagedFile := false
	for _, file := range d.Files {
		if file.IsStaged() {
			hasStagedFile = true
		} else {
			hasUnstagedFile = true
		}

		if hasStagedFile && hasUnstagedFile {
			return PartiallyStaged
		}
	}

	if hasStagedFile && !hasUnstagedFile {
		return FullyStaged
	}

	return Unstaged
}

func (d Directory) IsDirectory() bool { return true }

func (d Directory) GetFilePath() string { return d.Filepath }
