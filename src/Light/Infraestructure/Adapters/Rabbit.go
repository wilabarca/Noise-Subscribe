package adapters

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// Nueva función para inicializar el RabbitMQAdapter
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

// Método para consumir mensajes de RabbitMQ
func (r *RabbitMQAdapter) Consume() (<-chan amqp.Delivery, error) {
	queueName := "sensor_light"
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

	// Consumir los mensajes de la cola
	messages, err := r.channel.Consume(
		queueName, // nombre de la cola
		"",        // consumer tag
		true,      // auto ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("❌ Error al consumir los mensajes de la cola:", err)
		return nil, err
	}

	log.Printf("✅ Consumiendo mensajes de la cola '%s'\n", queueName)
	return messages, nil
}

// Método para cerrar la conexión y el canal de RabbitMQ
func (r *RabbitMQAdapter) Close() {
	r.channel.Close()
	r.connection.Close()
}
