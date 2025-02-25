package publisher

import (
	"github.com/marianozunino/go-publish/internal/models"

	"github.com/streadway/amqp"
)

// PublishMessage publishes a single message to the queue
func PublishMessage(ch *amqp.Channel, queueName string, msg models.RawMessage) error {
	return ch.Publish(
		"",        // exchange (empty for direct to queue)
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  msg.Properties.ContentType,
			DeliveryMode: uint8(msg.Properties.DeliveryMode),
			Priority:     uint8(msg.Properties.Priority),
			Body:         []byte(msg.Payload),
		})
}
