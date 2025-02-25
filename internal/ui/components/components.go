package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// RenderQueueInfo creates the queue information box
func RenderQueueInfo(queueName string, insecureTLS bool, styles map[string]lipgloss.Style) string {
	queueInfo := styles["subtitle"].Render("Queue Connection")
	queueInfo += "\n" + styles["info"].Render(queueName)

	if insecureTLS {
		queueInfo += "\n" + styles["warning"].Render("⚠️  Insecure mode: TLS certificate validation disabled")
	}

	return styles["box"].Render(queueInfo)
}

// RenderProgressSection creates the progress section
func RenderProgressSection(
	current int,
	total int,
	progressBar string,
	styles map[string]lipgloss.Style,
) string {
	progressSection := styles["subtitle"].Render("Progress")

	// Add the progress bar
	progressSection += "\n" + progressBar

	// Progress percentage & counts
	progress := float64(current) / float64(total)
	progressDetails := fmt.Sprintf("%d/%d messages (%.1f%%)",
		current, total, progress*100)
	progressSection += "\n" + styles["info"].Render(progressDetails)

	return styles["box"].Render(progressSection)
}

// RenderStatsSection creates the statistics section
func RenderStatsSection(
	isPaused bool,
	successCount int,
	errorCount int,
	delay time.Duration,
	msgPerSec float64,
	eta time.Duration,
	elapsed time.Duration,
	contentWidth int,
	styles map[string]lipgloss.Style,
) string {
	statsSection := styles["subtitle"].Render("Statistics")

	// Status line with emoji
	var statusLine string
	if isPaused {
		statusLine = styles["paused"].Render("⏸️  PAUSED")
	} else {
		statusLine = styles["running"].Render("▶️  RUNNING")
	}

	// Success and error counts
	successLine := fmt.Sprintf("✅ Success: %d", successCount)
	errorLine := fmt.Sprintf("❌ Errors: %d", errorCount)

	// ETA line (only show meaningful ETA when not paused)
	etaLine := "⏰ ETA: "
	if isPaused {
		etaLine += "Paused"
	} else {
		etaLine += eta.Round(time.Second).String()
	}

	// Create a grid layout for stats
	statsGrid := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(contentWidth/2-4).Render(
			statusLine+"\n"+
				styles["success"].Render(successLine)+"\n"+
				(func() string {
					if errorCount > 0 {
						return styles["error"].Render(errorLine)
					}
					return styles["info"].Render(errorLine)
				})()),
		lipgloss.NewStyle().Width(contentWidth/2-4).Render(
			styles["info"].Render(fmt.Sprintf("⏱️  Delay: %s", delay.Round(time.Millisecond)))+"\n"+
				styles["info"].Render(fmt.Sprintf("🚀 Speed: %.1f msg/sec", msgPerSec))+"\n"+
				styles["info"].Render(etaLine)),
	)

	statsSection += "\n" + statsGrid

	timeInfo := styles["info"].Render(fmt.Sprintf("Elapsed time: %s", elapsed.Round(time.Second)))
	statsSection += "\n" + timeInfo

	return styles["box"].Render(statsSection)
}

// RenderErrorBox creates an error box if there's an error
func RenderErrorBox(lastError string, styles map[string]lipgloss.Style) string {
	if lastError == "" {
		return ""
	}

	errorBox := styles["subtitle"].Render("Last Error")
	errorBox += "\n" + styles["error"].Render(lastError)
	return styles["box"].Copy().BorderForeground(styles["error"].GetForeground()).Render(errorBox)
}

// RenderControlsBox creates the controls box
func RenderControlsBox(styles map[string]lipgloss.Style) string {
	controlsText := styles["subtitle"].Render("Controls") + "\n" +
		styles["dimmed"].Render("SPACE") + " Pause/Resume | " +
		styles["dimmed"].Render("+") + " Increase Delay | " +
		styles["dimmed"].Render("-") + " Decrease Delay | " +
		styles["dimmed"].Render("q") + " Quit"

	return styles["controlsBox"].Render(controlsText)
}

// RenderEmptyState creates a message for when there are no messages to process
func RenderEmptyState(styles map[string]lipgloss.Style) string {
	return styles["box"].Render("No messages to process.") + "\n\n" +
		styles["dimmed"].Render("Press q to quit.")
}
