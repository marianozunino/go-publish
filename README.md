# GoPublish

A modern interactive CLI tool for publishing messages to RabbitMQ queues with style.

![GoPublish Screenshot](https://i.imgur.com/example-screenshot.png)

## Features

- **Interactive UI**: Real-time progress visualization with a stylish terminal interface
- **Controlled Publishing**: Adjust message publishing speed on-the-fly
- **Secure Connections**: Support for both AMQP and AMQPS connections
- **Statistics Monitoring**: View real-time throughput, success rates and estimated completion time
- **Pause & Resume**: Control message flow during publishing
- **TLS Support**: Connect to secure brokers (with optional certificate validation skipping)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/gopublish.git
cd gopublish

# Build the binary
go build -o gopublish

# Optional: Install to your $GOPATH/bin
go install
```

### With Go Install

```bash
go install github.com/yourusername/gopublish@latest
```

## Usage

### Basic Usage

```bash
# Publish messages from a file to a queue
gopublish --input messages.txt --queue your-queue-name
```

### Command-Line Options

```
Usage:
  gopublish [flags]

Flags:
  -d, --dry-run           Process file but don't send messages
  -h, --help              Help for gopublish
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
gopublish -i messages.json -q important-messages
```

### Connect to Remote RabbitMQ Server

```bash
gopublish -u amqp://user:password@rabbitmq-server:5672/
```

### Connect to AMQPS with Self-Signed Certificate

```bash
gopublish -u amqps://user:password@secure-rabbitmq:5671/ -k
```

### Dry Run (Test without Publishing)

```bash
gopublish -i messages.json --dry-run
```

## Project Structure

```
gopublish/
├── cmd/
│   └── root.go           # Cobra command definitions and CLI setup
├── internal/
│   ├── models/
│   │   └── message.go    # Message data models and file parsing
│   ├── publisher/
│   │   └── publisher.go  # Message publishing logic
│   └── ui/
│       └── tui.go        # Terminal UI implementation
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
