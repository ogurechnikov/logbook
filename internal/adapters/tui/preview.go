package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PreviewModel struct {
	viewport viewport.Model
	content  string
}

func NewPreviewModel() PreviewModel {
	vp := viewport.New(0, 0)
	return PreviewModel{
		viewport: vp,
	}
}

func (m PreviewModel) Init() tea.Cmd {
	return nil
}

func (m PreviewModel) Update(msg tea.Msg) (PreviewModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width / 2
		m.viewport.Height = msg.Height - 2
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *PreviewModel) SetContent(content string) {
	m.content = content
	m.viewport.SetContent(content)
}

func (m PreviewModel) View() string {
	if m.content == "" {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("Выберите заметку для просмотра")
	}

	return m.viewport.View()
}
