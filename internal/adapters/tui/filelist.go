package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileListModel struct {
	files  []string
	cursor int
	width  int
	height int
}

func NewFileListModel(files []string) FileListModel {
	return FileListModel{
		files:  files,
		cursor: 0,
	}
}

func (m FileListModel) Init() tea.Cmd {
	return nil
}

func (m FileListModel) Update(msg tea.Msg) (FileListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height - 2

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m FileListModel) View() string {
	if len(m.files) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("Нет заметок. Нажми 'n' чтобы создать.")
	}

	var b strings.Builder
	for i, file := range m.files {
		if i == m.cursor {
			b.WriteString("▸ ")
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffff")).
				Render(file))
		} else {
			b.WriteString("  ")
			b.WriteString(file)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m FileListModel) SelectedFile() string {
	if len(m.files) == 0 || m.cursor >= len(m.files) {
		return ""
	}
	return m.files[m.cursor]
}
