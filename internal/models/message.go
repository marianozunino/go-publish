package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// RawMessage represents the input file message format
type RawMessage struct {
	PayloadBytes    int               `json:"payload_bytes"`
	Redelivered     bool              `json:"redelivered"`
	Exchange        string            `json:"exchange"`
	RoutingKey      string            `json:"routing_key"`
	MessageCount    int               `json:"message_count"`
	Properties      MessageProperties `json:"properties"`
	Payload         string            `json:"payload"`
	PayloadEncoding string            `json:"payload_encoding"`
}

// MessageProperties represents AMQP message properties
type MessageProperties struct {
	Priority     int    `json:"priority"`
	DeliveryMode int    `json:"delivery_mode"`
	ContentType  string `json:"content_type"`
}

// ParseMessageFile reads and parses messages from a file
func ParseMessageFile(filePath string) ([]RawMessage, error) {
	// Open and read the input file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var messages []RawMessage

	// Parse each line as a JSON object
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		var msg RawMessage
		err := json.Unmarshal([]byte(line), &msg)
		if err != nil {
			return nil, fmt.Errorf("error parsing message at line %d: %w", lineNum, err)
		}
		messages = append(messages, msg)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("no valid messages found in file")
	}

	return messages, nil
}
