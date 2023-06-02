package pgAdapter

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	AccountUUID uuid.UUID
	Money       float64
	CreateTime  time.Time
	UpdateTime  time.Time
}

type Request struct {
	RequestID     uuid.UUID
	AccountUUID   uuid.UUID
	MoneyTransfer float64
	Status        string
	CreateTime    time.Time
	UpdateTime    time.Time
}
