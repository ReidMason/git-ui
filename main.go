package main

import (
	"fmt"
	"git-ui/internal/git"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Align(lipgloss.Position(0.5))

	greyOutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4b5161"))
)

type model struct {
	ldiff     []git.DiffLine
	rdiff     []git.DiffLine
	lviewport viewport.Model
	rviewport viewport.Model
	width     int
	ready     bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filepath")
		return
	}

	filepath := os.Args[1]

	p := tea.NewProgram(
		initialModel(filepath),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func styleLine(line git.DiffLine, width int) string {
	lineString := line.Content[:min(width-7, len(line.Content))]

	additionStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#3f534f")).
		Width(width)

	removalStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#6f2e2d")).
		Width(width)

	blankStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#31343b")).
		Width(width)

	if line.Type == git.Removal {
		lineString = removalStyle.Render(lineString)
	} else if line.Type == git.Addition {
		lineString = additionStyle.Render(lineString)
	} else if line.Type == git.Blank {
		lineString = blankStyle.Render(lineString)
	}

	return lineString
}

func styleDiff(diff []git.DiffLine, width int) string {
	diffString := ""
	count := 1

	numberOfLines := 0
	for _, line := range diff {
		if line.Type != git.Blank {
			numberOfLines++
		}
	}
	lineNumberPadding := len(fmt.Sprint(numberOfLines))

	for _, line := range diff {
		// TODO: This line number padding is probably really slow it could do with improving
		isBlank := line.Type == git.Blank
		lineNumber := strings.Repeat(" ", lineNumberPadding)
		if !isBlank {
			lengthOfCurrentNumber := len(fmt.Sprint(count))
			lineNumber = strings.Repeat(" ", lineNumberPadding-lengthOfCurrentNumber)
			lineNumber += fmt.Sprint(count)
			count++
		}
		lineNumber += "│"

		diffString += greyOutStyle.Render(lineNumber) + styleLine(line, width) + "\n"
	}

	return diffString
}

func initialModel(filepath string) model {
	rawDiff := git.GetRawDiff(filepath)
	diff := git.GetDiff(rawDiff)

	return model{
		ldiff: diff.Diff1,
		rdiff: diff.Diff2,
		ready: false,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		offset := 2
		m.width = msg.Width
		width := m.width/2 - offset
		lineWidth := width - columnStyle.GetHorizontalPadding()
		height := msg.Height - columnStyle.GetVerticalPadding() - 5

		if !m.ready {
			m.lviewport = viewport.New(width, height)
			m.lviewport.YPosition = 10
			ldiff := styleDiff(m.ldiff, lineWidth)
			m.lviewport.SetContent(ldiff)

			m.rviewport = viewport.New(width, height)
			m.rviewport.YPosition = 10

			rdiff := styleDiff(m.rdiff, lineWidth)
			m.rviewport.SetContent(rdiff)

			columnStyle.Width(width)
			m.ready = true
		} else {
			columnStyle.Width(width)
			columnStyle.Height(height)

			ldiff := styleDiff(m.ldiff, lineWidth)
			m.lviewport.SetContent(ldiff)

			rdiff := styleDiff(m.rdiff, lineWidth)
			m.rviewport.SetContent(rdiff)

			m.lviewport.Width = width
			m.lviewport.Height = height

			m.rviewport.Width = width
			m.rviewport.Height = height
		}
	}

	m.lviewport, cmd = m.lviewport.Update(msg)
	cmds = append(cmds, cmd)

	m.rviewport, cmd = m.rviewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	headerStlying := headerStyle.Width(m.width - 2)
	header := headerStlying.Render("Git diff")

	leftView := columnStyle.Render(m.lviewport.View())
	rightView := columnStyle.Render(m.rviewport.View())

	mainBody := lipgloss.JoinHorizontal(lipgloss.Left, leftView, rightView)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainBody)
}
