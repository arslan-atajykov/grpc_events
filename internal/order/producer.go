package order

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &Producer{writer: w}
}

func (p *Producer) PublishOrder(ctx context.Context, o *Order) error {
	payload, err := json.Marshal(o)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(o.Customer),
		Value: payload,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}

	log.Printf("Published order event to Kafka: %v", o)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
