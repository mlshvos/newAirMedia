package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"newAirMedia/internal/A/serverHandles"
	"newAirMedia/internal/util/kafkaAdapter"
	"newAirMedia/internal/util/pgAdapter"
	"strconv"
)

func parseCli(port *int) {
	flag.IntVar(port, "port", 8080, "port to listen on")

	flag.Parse()

	log.Printf("[Flags] port is %v", *port)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			pgAdapter.Exit()
			kafkaAdapter.ExitProducer()
			log.Fatal(err)
		}
	}()

	var port int
	parseCli(&port)

	kafkaAdapter.InitProducer()
	pgAdapter.Init()

	router := httprouter.New()
	router.GET("/ping", serverHandles.Ping) // liveness probe
	router.POST("/transfer-money/:accountUuid/:requestId", serverHandles.TransferMoney)

	router.PanicHandler = func(res http.ResponseWriter, req *http.Request, err interface{}) {
		log.Printf("[F] %s\n", err)
		res.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("[+] A successfully initialized, listening %v port", port)

	listenPort := ":" + strconv.Itoa(port)
	if err := http.ListenAndServe(listenPort, serverHandles.Limit(router)); err != nil {
		log.Println("[Critical] ListenAndServe error: ", err)
	}

	pgAdapter.Exit()
	kafkaAdapter.ExitProducer()
	return
}
