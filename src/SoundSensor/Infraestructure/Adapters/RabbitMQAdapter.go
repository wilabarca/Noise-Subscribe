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

func (r *RabbitMQAdapter) Consume() (<-chan amqp.Delivery, error) {
	queueName := "sensor_sound"
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  
		false, 
		false, 
		false, 
		nil,
	)
	if err != nil {
		log.Println("❌ Error al declarar la cola:", err)
		return nil, err
	}

	messages, err := r.channel.Consume(
		"sensor_sound",
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println("❌ Error al consumir los mensajes de la cola:", err)
		return nil, err
	}

	log.Printf("✅ Consumiendo mensajes de la cola '%s'\n", queueName)
	return messages, nil
}

func (r *RabbitMQAdapter) Close() {
	r.channel.Close()
	r.connection.Close()
}
