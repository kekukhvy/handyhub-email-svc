package queue

import (
	"encoding/json"
	"fmt"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var log = logrus.StandardLogger()

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *config.RabbitMQConfig
}

func NewRabbitMQ(cfg *config.RabbitMQConfig) (*RabbitMQ, error) {
	log.Info("Connecting to RabbitMQ...")
	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		log.WithError(err).Errorf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.WithError(err).Errorf("Failed to open a channel: %v", err)
		return nil, err
	}

	log.Infof("Connected to RabbitMQ at %s", cfg.Url)

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		cfg:     cfg,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.WithError(err).Error("Failed to close RabbitMQ channel")
			return err
		} else {
			log.Info("RabbitMQ channel closed")
			return nil
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.WithError(err).Error("Failed to close RabbitMQ connection")
			return err
		} else {
			log.Info("RabbitMQ connection closed")
			return nil
		}
	}

	return nil
}

func (r *RabbitMQ) SetupQueue() error {
	err := r.channel.ExchangeDeclare(
		r.cfg.Exchange,
		r.cfg.ExchangeType,
		r.cfg.Durable,
		r.cfg.AutoDelete,
		r.cfg.Internal,
		r.cfg.NoWait,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = r.channel.QueueDeclare(
		r.cfg.EmailQueue,
		r.cfg.Durable,
		r.cfg.AutoDelete,
		r.cfg.Exclusive,
		r.cfg.NoWait,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}
	err = r.channel.QueueBind(
		r.cfg.EmailQueue,
		r.cfg.RoutingKey,
		r.cfg.Exchange,
		r.cfg.NoWait,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	log.Infof("Queue %s declared and bound to exchange %s with routing key %s", r.cfg.EmailQueue, r.cfg.Exchange, r.cfg.RoutingKey)
	return nil
}

func (r *RabbitMQ) ConsumeMessages() (<-chan amqp.Delivery, error) {
	err := r.channel.Qos(
		r.cfg.PrefetchCount,
		r.cfg.PrefetchSize,
		r.cfg.Global,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %v", err)
	}

	msg, err := r.channel.Consume(
		r.cfg.EmailQueue,
		r.cfg.Consumer,
		r.cfg.AutoAck,
		r.cfg.Exclusive,
		r.cfg.NoLocal,
		r.cfg.NoWait,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %v", err)
	}

	log.Infof("Consuming messages from queue %s", r.cfg.EmailQueue)
	return msg, nil
}

func (r *RabbitMQ) ParseMessage(body []byte) (*models.QueueMessage, error) {
	var msg models.QueueMessage
	err := json.Unmarshal(body, &msg)
	if err != nil {
		log.WithError(err).Errorf("Failed to parse message: %v", err)
		return nil, err
	}

	return &msg, nil
}
