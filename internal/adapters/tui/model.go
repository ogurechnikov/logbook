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

	creating      bool
	editing       bool
	confirmDelete bool

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

		if m.creating || m.editing {
			m.textarea.SetWidth(msg.Width - 4)
			m.textarea.SetHeight(msg.Height - 4)
		}

		m.fileList, _ = m.fileList.Update(msg)
		m.preview, _ = m.preview.Update(msg)
		return m, nil

	case tea.KeyMsg:
		if m.confirmDelete {
			return m, m.handleDeleteConfirm(msg)
		}

		if m.creating || m.editing {
			if msg.String() == "ctrl+s" || msg.String() == "esc" {
				return m, m.handleEditMode(msg)
			}
			m.textarea, _ = m.textarea.Update(msg)
			return m, nil
		}

		// Режим просмотра
		handled, cmd := m.handleViewMode(msg)
		if handled {
			return m, cmd
		}

		// Если не обработано — делегируем в fileList (j/k)
		m.fileList, cmd = m.fileList.Update(msg)
		if file := m.fileList.SelectedFile(); file != "" {
			content, err := m.fileRepo.Read(m.vaultPath, file)
			if err == nil {
				m.preview.SetContent(content)
			}
		}
		return m, cmd
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

	if m.confirmDelete {
		file := m.fileList.SelectedFile()
		dialog := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render(fmt.Sprintf("Удалить '%s'? (y/N)", file))

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			dialog,
		)
	}

	if m.creating || m.editing {
		headerText := "⚓ Новая заметка"
		if m.editing {
			headerText = "⚓ Редактирование"
		}
		header := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Render(headerText)

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
		Render("j/k: навигация · n: новая · e: ред. · d: удалить · q: выход")

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
	contentHeight := m.height - 2

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

// handleViewMode обрабатывает клавиши в режиме просмотра.
// Возвращает true если клавиша была обработана, и команду.
func (m *Model) handleViewMode(msg tea.KeyMsg) (bool, tea.Cmd) {
	switch msg.String() {
	case "n":
		m.startCreate()
		return true, textarea.Blink
	case "e":
		m.startEdit()
		return true, textarea.Blink
	case "d":
		m.confirmDelete = true
		return true, nil
	case "q", "ctrl+c":
		return true, tea.Quit
	}
	return false, nil
}

// handleEditMode обрабатывает клавиши в режиме создания/редактирования
func (m *Model) handleEditMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+s":
		m.saveNote()
	case "esc":
		m.cancelEdit()
	}
	return nil
}

// handleDeleteConfirm обрабатывает клавиши в режиме подтверждения удаления
func (m *Model) handleDeleteConfirm(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "y":
		m.deleteSelectedFile()
	case "n", "esc":
		m.confirmDelete = false
	}
	return nil
}

func (m *Model) startCreate() {
	m.creating = true
	m.textarea.Reset()
	m.textarea.Focus()
	m.textarea.SetWidth(m.width - 4)
	m.textarea.SetHeight(m.height - 4)
}

func (m *Model) startEdit() {
	file := m.fileList.SelectedFile()
	if file == "" {
		return
	}

	content, err := m.fileRepo.Read(m.vaultPath, file)
	if err != nil {
		return
	}

	m.editing = true
	m.textarea.Reset()
	m.textarea.SetValue(file + "\n" + content)
	m.textarea.Focus()
	m.textarea.SetWidth(m.width - 4)
	m.textarea.SetHeight(m.height - 4)
}

func (m *Model) saveNote() {
	content := m.textarea.Value()
	if content == "" {
		m.creating = false
		m.editing = false
		return
	}

	lines := strings.Split(content, "\n")
	filename := strings.TrimSpace(lines[0])
	if filename == "" {
		filename = "новая-заметка"
	}
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	body := ""
	if len(lines) > 1 {
		body = strings.Join(lines[1:], "\n")
	}

	m.fileRepo.Save(m.vaultPath, filename, body)
	m.refreshFileList()
	m.creating = false
	m.editing = false
}

func (m *Model) cancelEdit() {
	m.creating = false
	m.editing = false
	m.refreshFileList()
}

func (m *Model) deleteSelectedFile() {
	file := m.fileList.SelectedFile()
	if file == "" {
		return
	}

	m.fileRepo.Delete(m.vaultPath, file)
	m.refreshFileList()
	m.confirmDelete = false
}

func (m *Model) refreshFileList() {
	files, _ := m.fileRepo.List(m.vaultPath)
	m.fileList = NewFileListModel(files)
}
