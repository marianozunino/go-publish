package cmd

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/marianozunino/go-publish/internal/models"
	"github.com/marianozunino/go-publish/internal/ui"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

var (
	// Flags
	inputFile     string
	queueName     string
	amqpURI       string
	dryRun        bool
	initialDelay  int
	skipTLSVerify bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-publish",
	Short: "Resend messages to a RabbitMQ queue with interactive controls",
	Long: `A flexible utility for resending messages to a RabbitMQ queue with
interactive controls for speed and pausing.

Controls in the UI:
  SPACE - Pause/Resume publishing
  +     - Increase delay between messages
  -     - Decrease delay between messages
  q     - Quit the application`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse messages from the input file
		messages, err := models.ParseMessageFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully parsed %d messages from %s\n", len(messages), inputFile)

		if dryRun {
			fmt.Printf("Dry run mode. %d messages would have been sent to queue: %s\n",
				len(messages), queueName)
			return
		}

		// Connect to RabbitMQ
		conn, ch, err := connectToRabbitMQ(amqpURI, queueName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()
		defer ch.Close()

		// Start the interactive UI
		if err := ui.StartTUI(messages, conn, ch, queueName, initialDelay, skipTLSVerify); err != nil {
			fmt.Fprintf(os.Stderr, "UI Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewUpdateCmd())
	return rootCmd.Execute()
}

func init() {
	// Define flags
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "paste.txt",
		"Input file containing messages")
	rootCmd.PersistentFlags().StringVarP(&queueName, "queue", "q", "member-dossier",
		"Target queue name")
	rootCmd.PersistentFlags().StringVarP(&amqpURI, "uri", "u",
		"amqp://guest:guest@localhost:5672/", "AMQP URI")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false,
		"Process file but don't send messages")
	rootCmd.PersistentFlags().IntVarP(&initialDelay, "delay", "m", 10,
		"Initial delay between messages in milliseconds")
	rootCmd.PersistentFlags().BoolVarP(&skipTLSVerify, "insecure", "k", false,
		"Skip TLS certificate validation for AMQPS connections")
}

// Helper function to connect to RabbitMQ and set up a channel
func connectToRabbitMQ(uri, queueName string) (*amqp.Connection, *amqp.Channel, error) {
	var conn *amqp.Connection
	var err error

	// Check if we need custom TLS config (for amqps:// URLs)
	if skipTLSVerify && strings.HasPrefix(uri, "amqps://") {
		// Create custom AMQP config with TLS settings
		cfg := amqp.Config{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		// Connect with the custom config
		conn, err = amqp.DialConfig(uri, cfg)
	} else {
		// Use standard connection
		conn, err = amqp.Dial(uri)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Ensure the target queue exists
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return conn, ch, nil
}
