package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
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

	creating bool
	textarea textarea.Model
}

func NewModel(vaultPath string, fileRepo domain.FileRepository) *Model {
	files, err := fileRepo.List(vaultPath)

	ta := textarea.New()
	ta.Placeholder = "Название заметки (первая строка)"
	ta.Focus()

	if err != nil {
		files = []string{}
	}
	return &Model{
		vaultPath: vaultPath,
		fileRepo:  fileRepo,
		fileList:  NewFileListModel(files),
		preview:   NewPreviewModel(),
		textarea:  ta,
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

		if m.creating {
			m.textarea.SetWidth(msg.Width - 4)
			m.textarea.SetHeight(msg.Height - 4)
		}

		m.fileList, _ = m.fileList.Update(msg)
		m.preview, _ = m.preview.Update(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "n":
			if !m.creating {
				m.creating = true
				m.textarea.Reset()
				m.textarea.Focus()
				m.textarea.SetWidth(m.width - 4)
				m.textarea.SetHeight(m.height - 4)
				return m, textarea.Blink
			}

		case "ctrl+s":
			if m.creating {
				content := m.textarea.Value()

				if content == "" {
					m.creating = false
					return m, nil
				}

				lines := strings.Split(content, "\n")
				filename := strings.TrimSpace(lines[0])
				if filename == "" {
					filename = "Новая заметка"
				}
				if !strings.HasSuffix(filename, ".md") {
					filename += ".md"
				}

				body := ""
				if len(lines) > 1 {
					body = strings.Join(lines[1:], "\n")
				}

				m.fileRepo.Save(m.vaultPath, filename, body)

				files, _ := m.fileRepo.List(m.vaultPath)
				m.fileList = NewFileListModel(files)
				m.creating = false
				return m, nil
			}
		case "esc":
			if m.creating {
				m.creating = false
				files, _ := m.fileRepo.List(m.vaultPath)
				m.fileList = NewFileListModel(files)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd

	if m.creating {
		m.textarea, cmd = m.textarea.Update(msg)
	} else {
		m.fileList, cmd = m.fileList.Update(msg)
		if file := m.fileList.SelectedFile(); file != "" {
			content, err := m.fileRepo.Read(m.vaultPath, file)
			if err == nil {
				m.preview.SetContent(content)
			}
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	if m.width == 0 {
		return "Загрузка..."
	}

	// Режим создания заметки
	if m.creating {
		header := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Render("⚓ Новая заметка")

		footer := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("Ctrl+S: сохранить · Esc: отменить")

		content := lipgloss.NewStyle().
			Padding(1).
			Render(m.textarea.View())

		contentHeight := m.height - 2
		content = lipgloss.NewStyle().
			Height(contentHeight).
			Render(content)

		return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
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
