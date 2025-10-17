package tui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/super-smooth/better-fg/internal/jobs"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	itemStyle = lipgloss.NewStyle()

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Padding(0, 1)
)

type item struct {
	job jobs.Job
}

func (i item) Title() string {
	emoji := getStatusEmoji(i.job.State)
	return fmt.Sprintf("[%d] %s %s", i.job.ID, i.job.Command, emoji)
}
func (i item) Description() string { return "" }

func getStatusEmoji(state string) string {
	switch state {
	case "suspended":
		return "⏸️"
	case "running":
		return "▶️"
	case "stopped", "Stopped":
		return "⏹️"
	case "done":
		return "✅"
	default:
		return "❓"
	}
}
func (i item) FilterValue() string { return i.job.Command }

type model struct {
	list     list.Model
	choice   *jobs.Job
	quitting bool
}

type customDelegate struct{}

func (d customDelegate) Height() int                             { return 1 }
func (d customDelegate) Spacing() int                            { return 0 }
func (d customDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := i.Title()

	if index == m.Index() {
		str = selectedItemStyle.Render(str)
	} else {
		str = itemStyle.Render(str)
	}

	fmt.Fprint(w, str)
}

func NewModel(jobs []jobs.Job) model {
	items := make([]list.Item, len(jobs))
	for i, job := range jobs {
		items[i] = item{job: job}
	}

	l := list.New(items, customDelegate{}, 0, 0)
	l.Title = "Select a job to bring to the foreground"
	l.Styles.Title = titleStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	return model{list: l}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = &i.job
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != nil {
		return ""
	}
	if m.quitting {
		return ""
	}
	return appStyle.Render(m.list.View())
}

func Run(jobs []jobs.Job) (*jobs.Job, error) {
	m := NewModel(jobs)
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	return finalModel.(model).choice, nil
}
