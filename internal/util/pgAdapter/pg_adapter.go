package pgAdapter

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"log"
)

var (
	DB *pg.DB
)

func Init() {
	Addr := fmt.Sprintf("%s:%d", host, port)
	DB = pg.Connect(&pg.Options{
		Addr:     Addr,
		User:     user,
		Password: password,
		Database: dbname,
	})
}

func Exit() {
	if err := DB.Close(); err != nil {
		log.Fatal(err)
	}
}
