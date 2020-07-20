package conn

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	port = 5432

	psqlUFOTrackerUsername = "postgres_ufotracker_user"
	psqlUFOTrackerPassword = "postgres_ufotracker_password"
	psqlUFOTrackerHost     = "postgres_ufotracker_host"
	psqlUFOTrackerSchema   = "postgres_ufotracker_dbname"
)

var (
	DB *sql.DB

	user     = os.Getenv(psqlUFOTrackerUsername)
	password = os.Getenv(psqlUFOTrackerPassword)
	host     = os.Getenv(psqlUFOTrackerHost)
	dbname   = os.Getenv(psqlUFOTrackerSchema)
)

func init() {
	var err error

	log.Println(host)
	log.Println(user)
	log.Println(dbname)

	info := fmt.Sprintf(`host=%s port=%d user=%s password=%s
		dbname=%s sslmode=disable`, host, port, user, password, dbname)
	DB, err = sql.Open("postgres", info)
	if err != nil {
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	log.Println(`Database successfully configured`)
}
