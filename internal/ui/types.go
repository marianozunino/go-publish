package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/streadway/amqp"
	"github.com/marianozunino/go-publish/internal/models"
)

// PublisherState holds the data relevant to the message publishing logic
type PublisherState struct {
	Messages      []models.RawMessage
	CurrentIndex  int
	TotalMessages int
	Delay         time.Duration
	Connection    *amqp.Connection
	Channel       *amqp.Channel
	QueueName     string
	InsecureTLS   bool
	LastError     string
}

// UIState holds the data relevant to the user interface
type UIState struct {
	Progress progress.Model
	IsPaused bool
	Width    int
	Height   int
	Theme    Theme
	Styles   Styles
}

// Statistics holds the data relevant to tracking performance and timing
type Statistics struct {
	SuccessCount    int
	ErrorCount      int
	StartTime       time.Time
	PauseStartTime  time.Time     // Track when pause starts
	TotalPausedTime time.Duration // Track total paused time
}

// Model represents the overall application state
type Model struct {
	Publisher PublisherState
	UI        UIState
	Stats     Statistics
}

// Define UI messages
type (
	tickMsg          time.Time
	publishResultMsg struct {
		success bool
		err     error
	}
)
