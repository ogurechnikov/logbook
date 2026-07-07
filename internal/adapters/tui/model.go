package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogurechnikov/logbook/internal/domain"
)

type Model struct {
	vaultPath string
	fileRepo  domain.FileRepository
	fileList  FileListModel
	preview   PreviewModel
	width     int
	height    int
}

func NewModel(vaultPath string, fileRepo domain.FileRepository) *Model {
	files, err := fileRepo.List(vaultPath)
	if err != nil {
		files = []string{}
	}
	return &Model{
		vaultPath: vaultPath,
		fileRepo:  fileRepo,
		fileList:  NewFileListModel(files),
		preview:   NewPreviewModel(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.fileList, _ = m.fileList.Update(msg)
		m.preview, _ = m.preview.Update(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.fileList, cmd = m.fileList.Update(msg)

	if file := m.fileList.SelectedFile(); file != "" {
		content, err := m.fileRepo.Read(m.vaultPath, file)
		if err == nil {
			m.preview.SetContent(content)
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	if m.width == 0 {
		return "Загрузка..."
	}

	// Верхняя панель
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Render(fmt.Sprintf("⚓ Logbook — %s", m.vaultPath))

	// Нижняя панель
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("j/k: навигация · n: новая · e: ред. · q: выход")

	// Разделитель
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#333333")).
		Render("│")

	// Панели с фиксированной шириной
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth - 1

	leftPanel := lipgloss.NewStyle().
		Width(leftWidth).
		Render(m.fileList.View())

	rightPanel := lipgloss.NewStyle().
		Width(rightWidth).
		Render(m.preview.View())

	// Основной контент: список + разделитель + предпросмотр
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Собираем: header, контент (растянут), footer (прижат книзу)
	contentHeight := m.height - 2 // минус header и footer

	mainContent = lipgloss.NewStyle().
		Height(contentHeight).
		Render(mainContent)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		mainContent,
		footer,
	)
}
