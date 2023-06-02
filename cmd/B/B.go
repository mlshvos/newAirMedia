package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"net/http"
	"newAirMedia/internal/A/serverHandles"
	"newAirMedia/internal/B/AHandler"
	"newAirMedia/internal/util/kafkaAdapter"
	"newAirMedia/internal/util/pgAdapter"
	"strconv"
	"sync"
	"time"
)

func parseCli(readersCount *int, requestHandlersMaxCount *int, port *int) {
	flag.IntVar(readersCount, "readersCount", 10, "kafka readers count")
	flag.IntVar(requestHandlersMaxCount, "requestHandlersMaxCount", 100, "request handlers max count")
	flag.IntVar(port, "port", 8081, "port to listen on")

	flag.Parse()

	log.Printf("[Flags] readersCount is %v", *readersCount)
	log.Printf("[Flags] requestHandlersMaxCount is %v", *requestHandlersMaxCount)
	log.Printf("[Flags] port is %v", *port)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			pgAdapter.Exit()
			kafkaAdapter.ExitReader()
			log.Fatal(err)
		}
	}()

	var (
		readersCount            int
		requestHandlersMaxCount int
		port                    int
	)
	parseCli(&readersCount, &requestHandlersMaxCount, &port)
	rand.Seed(time.Now().UnixNano())

	kafkaAdapter.InitReader()
	pgAdapter.Init()

	wg := &sync.WaitGroup{}

	router := httprouter.New()
	router.GET("/ping", serverHandles.Ping) // liveness probe

	listenPort := ":" + strconv.Itoa(port)
	go http.ListenAndServe(listenPort, router)

	reqSemaphore := make(chan int, requestHandlersMaxCount)
	for worker := 0; worker < readersCount; worker++ {
		wg.Add(1)
		go AHandler.Listen(wg, reqSemaphore)
	}
	wg.Wait()

	pgAdapter.Exit()
	kafkaAdapter.ExitReader()
	return
}
