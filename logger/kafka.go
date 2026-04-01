package logger

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter wraps a kafka-go writer for pushing logs
type KafkaWriter struct {
	writer *kafka.Writer
}

// NewKafkaWriter creates a Kafka writer connecting to the given brokers and writing to a topic
func NewKafkaWriter(brokers []string, topic string) *KafkaWriter {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        false, // Important: must be false to check for immediate write errors and trigger fallback
	}

	return &KafkaWriter{
		writer: w,
	}
}

// Write publishes a generic payload slice to Kafka
func (kw *KafkaWriter) Write(msg []byte) error {
	return kw.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: msg, // Only providing value, key is empty making it distribute over topic partitions
		},
	)
}

// Close gracefully closes the inner Kafka writer
func (kw *KafkaWriter) Close() error {
	return kw.writer.Close()
}
