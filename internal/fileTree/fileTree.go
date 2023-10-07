package filetree

type FileTreeItem interface {
	GetName() string
	IsExpanded() bool
	ToggleExpanded()
	Children() int
	GetStatus() string
	IsFullyStaged() bool
}

type FileTreeLine struct {
	Item  FileTreeItem
	Depth int
}

func New(item FileTreeItem, depth int) FileTreeLine {
	return FileTreeLine{Item: item, Depth: depth}
}

func Render(lines []string) string {
	output := ""
	activeLine := 0

	for i, line := range lines {
		if i == activeLine {
			output += "> "
		} else {
			output += "  "
		}
		output += line + "\n"
	}

	return output
}
