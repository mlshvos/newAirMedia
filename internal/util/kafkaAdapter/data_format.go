package kafkaAdapter

import "github.com/google/uuid"

type Message struct {
	Money       float64
	AccountUUID uuid.UUID
	RequestID   uuid.UUID
}
