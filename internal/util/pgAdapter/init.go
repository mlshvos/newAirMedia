package pgAdapter

import (
	"os"
	"strconv"
)

var (
	host     string
	port     int
	user     string
	password string
	dbname   string
)

func init() {
	host = os.Getenv("PG_HOST")
	var err error
	port, err = strconv.Atoi(os.Getenv("PG_PORT"))
	if err != nil {
		panic("Bad PG_PORT")
	}
	user = os.Getenv("PG_USER")
	password = os.Getenv("PG_PASSWORD")
	dbname = os.Getenv("PG_DBNAME")
}
