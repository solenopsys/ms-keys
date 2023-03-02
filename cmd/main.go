package main

import (
	_ "github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"ms-keys/internal"
	"os"
	"time"
)

var listenAddress = os.Getenv("LISTEN_ADDRESS")
var successRedirect = os.Getenv("SUCCESS_REDIRECT")
var errorRedirect = os.Getenv("ERROR_REDIRECT")
var dbPath = os.Getenv("DB_PATH")
var authServiceHost = os.Getenv("AUTH_SERVICE_HOST")
var mailServerHost = os.Getenv("MAIL_SERVER_HOST")
var mailServerPort = os.Getenv("MAIL_SERVER_PORT")
var fromAddress = os.Getenv("MAIL_FROM_ADDRESS")

func main() {
	sessions := cache.New(5*time.Minute, 10*time.Minute)
	db := internal.DriveDb{
		Path: dbPath,
	}
	mail := internal.MailServ{
		AuthHost: authServiceHost,
		From:     fromAddress,
		Host:     mailServerHost,
		Port:     mailServerPort,
	}
	restServer := internal.RestServer{
		Sessions:      sessions,
		Db:            &db,
		MailServer:    &mail,
		SuccessUrl:    successRedirect,
		ErrorUrl:      errorRedirect,
		ListenAddress: listenAddress,
	}

	db.Open()
	defer db.Close()
	restServer.Run()
}
