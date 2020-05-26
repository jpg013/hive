package logging

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaOut represents ioWriter for Kafka Transport
type KafkaOut struct {
	writer *kafka.Writer
}

// Write method takes byte array and writes to kafka writer
func (out *KafkaOut) Write(p []byte) (n int, err error) {
	err = out.writer.WriteMessages(context.Background(),
		kafka.Message{Value: p},
	)

	return n, err
}

// returns new kafka writer
func makeKafkaWriter(brokers []string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    "app_log",
		Balancer: &kafka.LeastBytes{},
	})
}

// NewKafkaTransport configures and returns a kafka transport
func NewKafkaTransport(conf *TransportConfig) Transport {
	writer := makeKafkaWriter(conf.KafkaHosts)
	out := &KafkaOut{
		writer: writer,
	}

	return Transport{
		Level: conf.Level,
		Out:   out,
	}
}
