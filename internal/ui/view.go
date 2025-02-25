package ui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/marianozunino/go-publish/internal/ui/components"
)

// View renders the UI
func (m Model) View() string {
	if m.IsEmpty() {
		return components.RenderEmptyState(m.getStylesMap())
	}

	// Set flexible width based on terminal size
	contentWidth := m.UI.Width - 4 // account for borders
	if contentWidth < 60 {
		contentWidth = 60
	}

	// Calculate elapsed time correctly, accounting for pauses
	elapsed := m.calculateElapsedTime()

	// Calculate progress
	progress := float64(m.Publisher.CurrentIndex) / float64(m.Publisher.TotalMessages)

	// Calculate messages per second
	msgPerSec := m.calculateMessageRate(elapsed)

	// Calculate estimated time remaining
	eta := m.calculateETA(msgPerSec)

	// Build the view
	s := m.UI.Styles.Title.Render("🐰 RabbitMQ Message Resender 🐰")

	// Queue information box
	s += components.RenderQueueInfo(
		m.Publisher.QueueName,
		m.Publisher.InsecureTLS,
		m.getStylesMap(),
	)

	// Progress section
	// Adapt progress bar to terminal width
	m.UI.Progress.Width = contentWidth - 8
	progressBar := m.UI.Progress.ViewAs(progress)

	s += components.RenderProgressSection(
		m.Publisher.CurrentIndex,
		m.Publisher.TotalMessages,
		progressBar,
		m.getStylesMap(),
	)

	// Stats box
	s += components.RenderStatsSection(
		m.UI.IsPaused,
		m.Stats.SuccessCount,
		m.Stats.ErrorCount,
		m.Publisher.Delay,
		msgPerSec,
		eta,
		elapsed,
		contentWidth,
		m.getStylesMap(),
	)

	// Last error if any
	s += components.RenderErrorBox(m.Publisher.LastError, m.getStylesMap())

	// Controls box at bottom
	s += components.RenderControlsBox(m.getStylesMap())

	return s
}

// calculateElapsedTime calculates the elapsed time accounting for pauses
func (m Model) calculateElapsedTime() time.Duration {
	var elapsed time.Duration
	if m.UI.IsPaused {
		// If currently paused, only count time up to when pause started
		elapsed = m.Stats.PauseStartTime.Sub(m.Stats.StartTime) - m.Stats.TotalPausedTime
	} else {
		// If running, subtract total paused time from the total elapsed time
		elapsed = time.Since(m.Stats.StartTime) - m.Stats.TotalPausedTime
	}
	return elapsed.Round(time.Second)
}

// calculateMessageRate calculates the messages per second rate
func (m Model) calculateMessageRate(elapsed time.Duration) float64 {
	var msgPerSec float64
	if elapsed > 0 && m.Publisher.CurrentIndex > 0 {
		msgPerSec = float64(m.Publisher.CurrentIndex) / elapsed.Seconds()
	}
	return msgPerSec
}

// calculateETA calculates the estimated time remaining
func (m Model) calculateETA(msgPerSec float64) time.Duration {
	var eta time.Duration
	if msgPerSec > 0 && m.Publisher.CurrentIndex < m.Publisher.TotalMessages {
		remainingMsgs := m.Publisher.TotalMessages - m.Publisher.CurrentIndex
		etaSeconds := float64(remainingMsgs) / msgPerSec
		eta = time.Duration(etaSeconds) * time.Second
	}
	return eta
}

// getStylesMap converts the Styles struct to a map for easier access in components
func (m Model) getStylesMap() map[string]lipgloss.Style {
	return map[string]lipgloss.Style{
		"title":       m.UI.Styles.Title,
		"subtitle":    m.UI.Styles.Subtitle,
		"info":        m.UI.Styles.Info,
		"dimmed":      m.UI.Styles.Dimmed,
		"success":     m.UI.Styles.Success,
		"error":       m.UI.Styles.Error,
		"warning":     m.UI.Styles.Warning,
		"paused":      m.UI.Styles.Paused,
		"running":     m.UI.Styles.Running,
		"box":         m.UI.Styles.Box,
		"controlsBox": m.UI.Styles.ControlsBox,
	}
}
