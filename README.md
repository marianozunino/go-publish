# GoPublish

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/marianozunino/go-publish.git
cd go-publish

# Build the binary
go build -o go-publish

# Optional: Install to your $GOPATH/bin
go install
```

### With Go Install

```bash
go install github.com/marianozunino/go-publish@latest
```

## Usage

### Basic Usage

```bash
# Publish messages from a file to a queue
go-publish --input messages.txt --queue your-queue-name
```

### Command-Line Options

```
Usage:
  go-publish [flags]

Flags:
  -d, --dry-run           Process file but don't send messages
  -h, --help              Help for go-publish
  -i, --input string      Input file containing messages (default "paste.txt")
  -k, --insecure          Skip TLS certificate validation for AMQPS connections
  -m, --delay int         Initial delay between messages in milliseconds (default 10)
  -q, --queue string      Target queue name (default "member-dossier")
  -u, --uri string        AMQP URI (default "amqp://guest:guest@localhost:5672/")
```

### Interactive Controls

Once running, you can use the following keyboard controls:

- **Space**: Pause/Resume message publishing
- **+**: Increase delay between messages (slows down)
- **-**: Decrease delay between messages (speeds up)
- **q**: Quit the application

## Input File Format

GoPublish expects a JSON file with one message per line, following this format:

```json
{"payload_bytes": 1088, "redelivered": true, "exchange": "exchange_name", "routing_key": "routing.key", "message_count": 81080, "properties": {"priority": 0, "delivery_mode": 2, "content_type": "application/json"}, "payload": "{\"your\":\"message\"}", "payload_encoding": "string"}
```

The important fields are:
- `properties`: AMQP message properties
- `payload`: The actual message content to publish

## Examples

### Connect to Local RabbitMQ with Custom Queue

```bash
go-publish -i messages.json -q important-messages
```

### Connect to Remote RabbitMQ Server

```bash
go-publish -u amqp://user:password@rabbitmq-server:5672/
```

### Connect to AMQPS with Self-Signed Certificate

```bash
go-publish -u amqps://user:password@secure-rabbitmq:5671/ -k
```

### Dry Run (Test without Publishing)

```bash
go-publish -i messages.json --dry-run
```

## Project Structure

```
go-publish/
├── cmd/
│   └── root.go           # Cobra command definitions and CLI setup
├── internal/
│   ├── models/
│   │   └── message.go    # Message data models and file parsing
│   ├── publisher/
│   │   └── publisher.go  # Message publishing logic
│   └── ui/               # Terminal UI implementation
├── main.go               # Entry point
└── go.mod                # Module dependencies
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - Command-line interface framework
- [BubbleTea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Streadway AMQP](https://github.com/streadway/amqp) - RabbitMQ client library

## License

MIT License
