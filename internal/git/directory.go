package git

import (
	"sort"

	"github.com/ReidMason/git-ui/internal/colours"
	filetree "github.com/ReidMason/git-ui/internal/fileTree"

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
	status := d.GetStagedStatus()
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

func (d *Directory) Sort() {
	sort.Slice(d.Files, func(i, j int) bool {
		return d.Files[i].Name < d.Files[j].Name
	})

	for _, subDirectory := range d.Directories {
		subDirectory.Sort()
	}
}

func (d Directory) GetStagedStatus() StagedStatus {
	hasStaged := false
	hasUnstaged := false

	for _, subDirectory := range d.Directories {
		subDirectoryStatus := subDirectory.GetStagedStatus()
		if subDirectoryStatus == PartiallyStaged {
			return PartiallyStaged
		}

		if subDirectoryStatus == FullyStaged {
			hasStaged = true
		} else {
			hasUnstaged = true
		}

		if hasStaged && hasUnstaged {
			return PartiallyStaged
		}
	}

	for _, file := range d.Files {
		if file.IsStaged() {
			hasStaged = true
		} else {
			hasUnstaged = true
		}

		if hasStaged && hasUnstaged {
			return PartiallyStaged
		}
	}

	if hasStaged && !hasUnstaged {
		return FullyStaged
	}

	return Unstaged
}

func (d Directory) IsDirectory() bool { return true }

func (d Directory) GetFilePath() string { return d.Filepath }
func (d Directory) GetFirstFilePath() string {
	if len(d.Directories) > 0 {
		return d.Directories[0].Filepath
	}

	if len(d.Files) > 0 {
		return d.Files[0].GetFilePath()
	}

	return ""
}
