package kafkaAdapter

import (
	"log"
	"os"
	"strconv"
)

var (
	topic         = "money-transfer"
	host          string
	port          int
	brokerAddress = "localhost:9092"
)

func init() {
	topic = os.Getenv("KAFKA_TOPIC")
	host = os.Getenv("KAFKA_HOST")
	var err error
	port, err = strconv.Atoi(os.Getenv("KAFKA_PORT"))
	if err != nil {
		log.Fatal("Bad KAFKA_PORT")
	}
}
