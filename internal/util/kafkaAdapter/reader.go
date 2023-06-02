package kafkaAdapter

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

var reader *kafka.Reader

func InitReader() {
	Addr := fmt.Sprintf("%s:%d", host, port)

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{Addr},
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
		GroupID:  "1",
	})
}

func FetchMessage(ctx context.Context) (*kafka.Message, error) {
	m, err := reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("message at offset %d: %v\n", m.Offset, string(m.Value))

	return &m, nil
}

func CommitMessage(ctx context.Context, msg *kafka.Message) {
	if err := reader.CommitMessages(ctx, *msg); err != nil {
		log.Printf("[F] Failed to commit messages: %v\n", err)
	}
}

func ExitReader() {
	if err := reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
