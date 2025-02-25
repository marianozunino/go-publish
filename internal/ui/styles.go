package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme contains all the color definitions for the application
type Theme struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Success   lipgloss.Color
	Neutral   lipgloss.Color
	Dimmed    lipgloss.Color
}

// DefaultTheme returns the default color theme
func DefaultTheme() Theme {
	return Theme{
		Primary:   lipgloss.Color("#25A065"),
		Secondary: lipgloss.Color("#2C8C99"),
		Warning:   lipgloss.Color("#FF8C42"),
		Error:     lipgloss.Color("#E84855"),
		Success:   lipgloss.Color("#2EC4B6"),
		Neutral:   lipgloss.Color("#FFFDF5"),
		Dimmed:    lipgloss.Color("#ABABAB"),
	}
}

// Styles contains all the styled components for the UI
type Styles struct {
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Info        lipgloss.Style
	Dimmed      lipgloss.Style
	Success     lipgloss.Style
	Error       lipgloss.Style
	Warning     lipgloss.Style
	Paused      lipgloss.Style
	Running     lipgloss.Style
	Box         lipgloss.Style
	ControlsBox lipgloss.Style
}

// DefaultStyles returns the default styles for the UI
func DefaultStyles(theme Theme) Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1).
			Align(lipgloss.Center),

		Subtitle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			MarginBottom(1),

		Info: lipgloss.NewStyle().
			Foreground(theme.Secondary),

		Dimmed: lipgloss.NewStyle().
			Foreground(theme.Dimmed),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),

		Paused: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),

		Running: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),

		Box: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Secondary).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1),

		ControlsBox: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Dimmed).
			MarginTop(1).
			Padding(1, 2).
			Align(lipgloss.Center),
	}
}

// AdjustStyles updates the styles based on the current terminal width
func AdjustStyles(s Styles, width int) Styles {
	contentWidth := width - 4 // account for borders
	if contentWidth < 60 {
		contentWidth = 60
	}

	s.Title = s.Title.Width(contentWidth)
	s.Box = s.Box.Width(contentWidth)
	s.ControlsBox = s.ControlsBox.Width(contentWidth)

	return s
}
