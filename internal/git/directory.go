package git

import filetree "git-ui/internal/fileTree"

type Directory struct {
	Name        string
	Parent      *Directory
	Files       []File
	Directories []*Directory
}

func newDirectory(name string, parent *Directory) *Directory {
	return &Directory{
		Name:        name,
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

func (d Directory) GetName() string   { return d.Name }
func (d Directory) IsDirectory() bool { return true }
func (d Directory) Children() int {
	count := len(d.Files)
	for _, directory := range d.Directories {
		count++
		count += directory.Children()
	}

	return count
}

func (d Directory) GetFilePath() string {
	if len(d.Files) > 0 {
		return d.Files[0].GetFilePath()
	}

	if len(d.Directories) > 0 {
		return d.Directories[0].GetFilePath()
	}

	return ""
}
