package serverHandles

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"newAirMedia/internal/util/kafkaAdapter"
	"newAirMedia/internal/util/pgAdapter"
	"strconv"
)

func Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, _ = fmt.Fprint(w, "Pong\n")
	log.Println("[+] Ponged!")
}

func TransferMoney(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	accountUuid, err := uuid.Parse(ps.ByName("accountUuid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[F] Failed to parse accountUuid: %v\n", err)
		return
	}

	requestId, err := uuid.Parse(ps.ByName("requestId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[F] Failed to parse requestId: %v\n", err)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[F] Failed read body: %v\n", err)
		return
	}

	transferMoney, err := strconv.ParseFloat(string(bytes), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[F] Bad money to transfer: %v\n", err)
		return
	}

	request := new(pgAdapter.Request)
	err = pgAdapter.DB.Model(request).Where("request_id = ?", requestId.String()).Select()
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[F] Request with id=%s already exist", requestId)
		return
	}

	msg := kafkaAdapter.Message{
		RequestID:   requestId,
		AccountUUID: accountUuid,
		Money:       transferMoney,
	}
	err = kafkaAdapter.Produce(&msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[F] Error while sending msg to kafka: %v\n", err)
		return
	}
}
