package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/streadway/amqp"
	"github.com/marianozunino/go-publish/internal/models"
)

// NewModel initializes the application model
func NewModel(messages []models.RawMessage, conn *amqp.Connection, ch *amqp.Channel, queueName string, delayMs int, insecureTLS bool) Model {
	// Configure progress bar with custom style
	theme := DefaultTheme()
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
		progress.WithoutPercentage(),
	)

	// Set progress bar colors
	p.FullColor = string(theme.Primary)
	p.EmptyColor = string(theme.Dimmed)

	return Model{
		Publisher: PublisherState{
			Messages:      messages,
			CurrentIndex:  0,
			TotalMessages: len(messages),
			Delay:         time.Duration(delayMs) * time.Millisecond,
			Connection:    conn,
			Channel:       ch,
			QueueName:     queueName,
			InsecureTLS:   insecureTLS,
			LastError:     "",
		},
		UI: UIState{
			Progress: p,
			IsPaused: false,
			Width:    80,
			Height:   24,
			Theme:    theme,
			Styles:   DefaultStyles(theme),
		},
		Stats: Statistics{
			SuccessCount:    0,
			ErrorCount:      0,
			StartTime:       time.Now(),
			PauseStartTime:  time.Time{},
			TotalPausedTime: 0,
		},
	}
}

// Init initializes the Bubble Tea program
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		publishMessageCmd(m),
	)
}

// IsComplete returns true if all messages have been processed
func (m Model) IsComplete() bool {
	return m.Publisher.CurrentIndex >= m.Publisher.TotalMessages
}

// IsEmpty returns true if there are no messages to process
func (m Model) IsEmpty() bool {
	return m.Publisher.TotalMessages == 0
}

// StartTUI initializes and runs the terminal UI
func StartTUI(messages []models.RawMessage, conn *amqp.Connection, ch *amqp.Channel, queueName string, delayMs int, insecureTLS bool) error {
	p := tea.NewProgram(
		NewModel(messages, conn, ch, queueName, delayMs, insecureTLS),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
