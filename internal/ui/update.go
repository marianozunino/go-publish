package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update processes user input and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	case tickMsg:
		return m, tickCmd()
	case publishResultMsg:
		return m.handlePublishResultMsg(msg)
	}

	return m, nil
}

// handleKeyMsg processes keyboard input
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case " ":
		return m.togglePause()
	case "+":
		return m.increaseDelay()
	case "-":
		return m.decreaseDelay()
	}
	return m, nil
}

// handleWindowSizeMsg updates the model when window size changes
func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.UI.Width = msg.Width
	m.UI.Height = msg.Height
	m.UI.Progress.Width = msg.Width - 20
	m.UI.Styles = AdjustStyles(m.UI.Styles, msg.Width)
	return m, nil
}

// handlePublishResultMsg processes the result of publishing a message
func (m Model) handlePublishResultMsg(msg publishResultMsg) (tea.Model, tea.Cmd) {
	if m.Publisher.CurrentIndex < m.Publisher.TotalMessages {
		m.Publisher.CurrentIndex++
		if msg.success {
			m.Stats.SuccessCount++
		} else {
			m.Stats.ErrorCount++
			if msg.err != nil {
				m.Publisher.LastError = msg.err.Error()
			}
		}
	}

	if m.IsComplete() {
		return m, nil
	}

	if !m.UI.IsPaused {
		return m, publishMessageCmd(m)
	}

	return m, nil
}

// togglePause toggles the pause state
func (m Model) togglePause() (tea.Model, tea.Cmd) {
	m.UI.IsPaused = !m.UI.IsPaused

	if m.UI.IsPaused {
		// When pausing, record the time pause started
		m.Stats.PauseStartTime = time.Now()
		return m, nil
	} else {
		// When resuming, add the paused duration to the total paused time
		m.Stats.TotalPausedTime += time.Since(m.Stats.PauseStartTime)
		return m, publishMessageCmd(m)
	}
}

// increaseDelay increases the delay between messages
func (m Model) increaseDelay() (tea.Model, tea.Cmd) {
	m.Publisher.Delay += 5 * time.Millisecond
	return m, nil
}

// decreaseDelay decreases the delay between messages
func (m Model) decreaseDelay() (tea.Model, tea.Cmd) {
	// Allow decreasing to 0ms
	if m.Publisher.Delay >= 5*time.Millisecond {
		m.Publisher.Delay -= 5 * time.Millisecond
	} else if m.Publisher.Delay > 0 {
		m.Publisher.Delay = 0
	}
	return m, nil
}
