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

func (r *RabbitMQAdapter) Consume(handler func(amqp.Delivery)) error {
	queueName := "sensor_air"

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
		return err
	}

	msgs, err := r.channel.Consume(
		queueName,
		"",    
		true,  
		false, 
		false, 
		false, 
		nil,
	)
	if err != nil {
		log.Println("❌ Error al consumir mensajes de la cola:", err)
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	log.Printf("✅ Escuchando mensajes en la cola '%s'\n", queueName)
	return nil
}

func (r *RabbitMQAdapter) Close() {
	r.channel.Close()
	r.connection.Close()
}
