package adapters

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQAdapter(amqpURL string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Println("❌ Error al conectar a RabbitMQ:", err)
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println("❌ Error al abrir un canal en RabbitMQ:", err)
		return nil, err
	}
	return &RabbitMQAdapter{
		connection: conn,
		channel:    ch,
	}, nil
}

func (r *RabbitMQAdapter) Publish(queueName string, body []byte) error {
	// Declarar la cola para asegurarnos de que existe
	_, err := r.channel.QueueDeclare(
		queueName,
		false, // no durable
		false, // auto-delete
		false, // no-exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Println("❌ Error al declarar la cola:", err)
		return err
	}

	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key = nombre de la cola
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("❌ Error al publicar en la cola", queueName, ":", err)
		return err
	}
	log.Printf("✅ Mensaje publicado en la cola '%s'\n", queueName)
	return nil
}

func (r *RabbitMQAdapter) Close() {
	r.channel.Close()
	r.connection.Close()
}
