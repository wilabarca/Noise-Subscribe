package adapters

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQAdapter struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (r *RabbitMQAdapter) Connect() any {
	panic("unimplemented")
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

func (r *RabbitMQAdapter) Consume(queueName string, handler func(amqp.Delivery)) error {
	// Declarar la cola antes de consumirla
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

	// Consumir los mensajes de la cola
	msgs, err := r.channel.Consume(
		queueName,
		"",    // consumer tag
		true,  // auto-ack (automáticamente confirmamos los mensajes)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Println("❌ Error al consumir de la cola:", err)
		return err
	}

	// Iniciar un goroutine para consumir los mensajes
	go func() {
		for msg := range msgs {
			// Aquí puedes procesar el mensaje que recibimos
			handler(msg)
		}
	}()

	return nil
}

func (r *RabbitMQAdapter) Close() {
	r.channel.Close()
	r.connection.Close()
}
