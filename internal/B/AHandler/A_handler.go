package AHandler

import (
	"context"
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/segmentio/kafka-go"
	"log"
	"math/rand"
	"newAirMedia/internal/util/kafkaAdapter"
	"newAirMedia/internal/util/pgAdapter"
	"sync"
	"time"
)

func waitForStatus() string {
	time.Sleep(30 * time.Second)
	if rand.Intn(2) == 1 {
		return "failed"
	} else {
		return "success"
	}
}

func processRequest(ctx context.Context, kafkaMsg *kafka.Message, msg *kafkaAdapter.Message, reqSemaphore chan int,
	request *pgAdapter.Request, account *pgAdapter.Account) {
	request.Status = waitForStatus()
	if request.Status == "failed" {
		log.Printf("[W] RequestId=%s failed", msg.RequestID)
		account.Money += msg.Money
		if err := pgAdapter.DB.RunInTransaction(context.TODO(), func(tx *pg.Tx) error {
			_, err := pgAdapter.DB.Model(request).
				Set("status = ?status").
				Where("request_id = ?request_id").
				Update()
			if err != nil {
				return err
			}
			_, err = pgAdapter.DB.Model(account).
				Set("money = ?money").
				Where("account_uuid = ?account_uuid").
				Update()
			return err
		}); err != nil {
			// Here is the most vulnerable place - should log it, hang alerts
			// and send to dead letter topic to update it manually.
			// possibly it make sense to create automatic service shutdowner
			// if such cases will be more than 5 for 5 minutes
			log.Printf("[Critical] Error while resetting account money: %v\n", err)
		}
	} else {
		log.Printf("[+] RequestId=%s successfully finished", msg.RequestID)
		_, err := pgAdapter.DB.Model(request).
			Set("status = ?status").
			Where("request_id = ?request_id").
			Update()
		if err != nil {
			// Here is the most vulnerable place - should log it, hang alerts
			// and send to dead letter topic to update it manually.
			// possibly it make sense to create automatic service shutdowner
			// if such cases will be more than 5 for 5 minutes
			log.Printf("[Critical] Error while updating request status: %v\n", err)
		}
	}
	kafkaAdapter.CommitMessage(ctx, kafkaMsg)
	<-reqSemaphore
}

func Handle(ctx context.Context, kafkaMsg *kafka.Message, msg *kafkaAdapter.Message, reqSemaphore chan int) {

	account := new(pgAdapter.Account)
	err := pgAdapter.DB.Model(account).Where("account_uuid = ?", msg.AccountUUID.String()).Select()
	if err != nil {
		log.Printf("[F] Error while selecting account from db: %v\n", err)
		return
	}

	if msg.Money > account.Money {
		log.Printf("[F] Account=%s reached his money limit", msg.AccountUUID)
		kafkaAdapter.CommitMessage(ctx, kafkaMsg)
		return
	}

	request := new(pgAdapter.Request)
	err = pgAdapter.DB.Model(request).Where("request_id = ?", msg.RequestID.String()).Select()
	if err == nil {
		log.Printf("[F] Something went wrong. Current request=%s is already exist in db\n",
			request.RequestID.String())
		return
	}

	request.RequestID = msg.RequestID
	request.AccountUUID = msg.AccountUUID
	request.MoneyTransfer = msg.Money

	account.Money -= msg.Money
	if err := pgAdapter.DB.RunInTransaction(context.TODO(), func(tx *pg.Tx) error {
		_, err = pgAdapter.DB.Model(request).Insert()
		if err != nil {
			return err
		}
		_, err := pgAdapter.DB.Model(account).
			Set("money = ?money").
			Where("account_uuid = ?account_uuid").
			Update()
		return err
	}); err != nil {
		log.Printf("[F] Error while creating request: %v\n", err)
		return
	}
	reqSemaphore <- 1
	go processRequest(ctx, kafkaMsg, msg, reqSemaphore, request, account)
}

func Listen(wg *sync.WaitGroup, reqSemaphore chan int) {
	defer wg.Done()

	ctx := context.Background()
	for {
		msg, err := kafkaAdapter.FetchMessage(ctx)
		if err != nil {
			break
		}

		msgJson := &kafkaAdapter.Message{}
		if err := json.Unmarshal(msg.Value, msgJson); err != nil {
			log.Printf("[F] Error while parsing msg from kafka: %v\n", err)
		} else {
			Handle(ctx, msg, msgJson, reqSemaphore)
		}
	}
}
