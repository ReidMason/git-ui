package git

type Directory struct {
	Name        string
	Parent      *Directory
	Files       []File
	Directories []*Directory
	Expanded    bool
}

func newDirectory(name string, parent *Directory) *Directory {
	return &Directory{
		Name:        name,
		Files:       make([]File, 0),
		Directories: make([]*Directory, 0),
		Expanded:    true,
		Parent:      parent,
	}
}

func (d Directory) GetName() string { return d.Name }
func (d Directory) Children() int {
	count := len(d.Files)
	for _, directory := range d.Directories {
		count++
		count += directory.Children()
	}

	return count
}
func (d Directory) GetStatus() string { return "" }
func (d Directory) IsExpanded() bool  { return d.Expanded }
func (d Directory) IsVisible() bool {
	if d.Parent == nil {
		return true
	}

	return d.Parent.IsVisible() && d.Parent.Expanded
}
func (d Directory) IsFullyStaged() bool {
	for _, file := range d.Files {
		if !file.IsFullyStaged() {
			return false
		}
	}

	for _, directory := range d.Directories {
		if !directory.IsFullyStaged() {
			return false
		}
	}

	return true
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

func (d *Directory) ToggleExpanded() { d.Expanded = !d.Expanded }
