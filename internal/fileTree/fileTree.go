package filetree

type FileTreeItem interface {
	GetName() string
	IsExpanded() bool
	ToggleExpanded()
	Children() int
	GetIndexStatus() rune
	GetWorkTreeStatus() rune
	IsFullyStaged() bool
}

type FileTreeLine struct {
	Item  FileTreeItem
	Depth int
}

func New(item FileTreeItem, depth int) FileTreeLine {
	return FileTreeLine{Item: item, Depth: depth}
}