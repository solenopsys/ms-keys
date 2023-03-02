package main

import (
	_ "github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"os"
	"time"
)

var listenAddress = os.Getenv("LISTEN_ADDRESS")
var successRedirect = os.Getenv("SUCCESS_REDIRECT")
var errorRedirect = os.Getenv("ERROR_REDIRECT")
var dbPath = os.Getenv("DB_PATH")

func sendMail(email string) {

}

var sessions = cache.New(5*time.Minute, 10*time.Minute)

func main() {
	db := DriveDb{path: dbPath}
	defer db.close()
	restServer := RestServer{
		sessions: sessions,
		db:       &db,
	}
	restServer.runServer()
}
