package kafkaAdapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

var producer *kafka.Conn

func InitProducer() {
	Addr := fmt.Sprintf("%s:%d", host, port)

	partition := 0

	var err error
	producer, err = kafka.DialLeader(context.Background(), "tcp", Addr, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader: ", err)
	}
}

func Produce(message *Message) error {
	msgBytes, err := json.Marshal(*message)
	if err != nil {
		return err
	}
	_, err = producer.WriteMessages(
		kafka.Message{Value: msgBytes},
	)
	return err
}

func ExitProducer() {
	if err := producer.Close(); err != nil {
		log.Fatal("failed to close producer:", err)
	}
}
