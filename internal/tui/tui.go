package tui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/super-smooth/better-fg/internal/jobs"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	itemStyle = lipgloss.NewStyle()

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true).
			Padding(0, 1)
)

type item struct {
	job jobs.Job
}

// Title returns the job command for display.
func (i item) Title() string {
	emoji := getStatusEmoji(i.job.State)
	return fmt.Sprintf("[%d] %s %s", i.job.ID, i.job.Command, emoji)
}

// Description returns empty string as we show everything in title.
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

// FilterValue returns the job command for filtering.
func (i item) FilterValue() string { return i.job.Command }

type model struct {
	list     list.Model
	choice   *jobs.Job
	quitting bool
}

type customDelegate struct{}

// Height returns the height of each list item.
func (d customDelegate) Height() int { return 1 }

// Spacing returns the spacing between list items.
func (d customDelegate) Spacing() int { return 0 }

// Update handles delegate updates.
func (d customDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render renders a list item.
//
//nolint:gocritic // Bubble Tea requires this signature
func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := i.Title()

	if index == m.Index() {
		str = selectedItemStyle.Render("▶ " + str)
	} else {
		str = itemStyle.Render("  " + str)
	}

	_, _ = fmt.Fprint(w, str)
}

// NewModel creates a new TUI model with the given jobs.
func NewModel(jobList []jobs.Job) model {
	items := make([]list.Item, len(jobList))
	for i, job := range jobList {
		items[i] = item{job: job}
	}

	l := list.New(items, customDelegate{}, 0, 0)
	l.Title = "Select a job to bring to the foreground"
	l.Styles.Title = titleStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	return model{list: l}
}

// Init initializes the model.
//
//nolint:gocritic // Bubble Tea requires this signature
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
//
//nolint:gocritic // Bubble Tea requires this signature
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

// View renders the model.
//
//nolint:gocritic // Bubble Tea requires this signature
func (m model) View() string {
	if m.choice != nil {
		return ""
	}

	if m.quitting {
		return ""
	}

	return appStyle.Render(m.list.View())
}

// Run starts the TUI and returns the selected job.
func Run(jobList []jobs.Job) (*jobs.Job, error) {
	m := NewModel(jobList)

	// Use /dev/tty to ensure proper terminal access when run through shell functions
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// Fallback to stderr if /dev/tty is not available
		p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

		finalModel, err := p.Run()
		if err != nil {
			return nil, err
		}

		return finalModel.(model).choice, nil
	}
	defer tty.Close()

	// Force color output
	lipgloss.SetColorProfile(termenv.TrueColor)

	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	return finalModel.(model).choice, nil
}
