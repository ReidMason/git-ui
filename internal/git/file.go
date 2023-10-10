package git

import "path/filepath"

type File struct {
	Name           string
	Parent         *Directory
	Dirpath        string
	IndexStatus    rune
	WorkTreeStatus rune
}

func (f File) GetName() string {
	return f.Name
}

func (f File) ToggleExpanded()     {}
func (f File) IsExpanded() bool    { return true }
func (f File) IsVisible() bool     { return f.Parent.IsVisible() && f.Parent.Expanded }
func (f File) Children() int       { return 0 }
func (f File) IsFullyStaged() bool { return f.WorkTreeStatus == '.' }
func (f File) GetStatus() string   { return string(f.IndexStatus) + string(f.WorkTreeStatus) }
func (f File) GetFilePath() string { return filepath.Join(f.Dirpath, f.Name) }

func newFile(filePath string, indexStatus, workTreeStatus rune) File {
	dirpath, filename := filepath.Split(filePath)
	dirpath = filepath.Clean(dirpath)

	return File{Name: filename, Dirpath: dirpath, IndexStatus: indexStatus, WorkTreeStatus: workTreeStatus}
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

	// log.Println(dirpath)
	newDir := newDirectory(dirpath[0], directory)
	addFile(newDir, dirpath[1:], newFile)
	newDir.Parent = directory
	directory.Directories = append(directory.Directories, newDir)
}
