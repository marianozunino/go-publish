package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marianozunino/go-publish/internal/publisher"
)

// Command to tick for UI updates
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Command to publish the next message
func publishMessageCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		if m.UI.IsPaused || m.IsComplete() {
			return nil
		}

		// Add delay between messages (only if greater than 0)
		if m.Publisher.Delay > 0 {
			time.Sleep(m.Publisher.Delay)
		}

		currentIdx := m.Publisher.CurrentIndex
		msg := m.Publisher.Messages[currentIdx]
		err := publisher.PublishMessage(m.Publisher.Channel, m.Publisher.QueueName, msg)

		return publishResultMsg{
			success: err == nil,
			err:     err,
		}
	}
}

// Command to adjust delay up
func increaseDelayCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		// Increase delay by 5ms
		return nil
	}
}

// Command to adjust delay down
func decreaseDelayCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		// Decrease delay by 5ms
		return nil
	}
}

// Command to toggle pause state
func togglePauseCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		// Toggle pause state and handle timing
		if m.UI.IsPaused {
			// We're resuming - continue publishing
			return publishMessageCmd(m)()
		}
		return nil
	}
}

// Command batch for all regular UI updates
func regularUpdateCmds(m Model) tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if !m.UI.IsPaused && !m.IsComplete() {
		cmds = append(cmds, publishMessageCmd(m))
	}
	return tea.Batch(cmds...)
}
