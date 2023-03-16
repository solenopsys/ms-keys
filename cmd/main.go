package main

import (
	"flag"
	_ "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"ms-keys/internal"
	"os"
	"time"
)

var Mode string

const DEV_MODE = "dev"

func init() {
	flag.StringVar(&Mode, "mode", "", "a string var")
}

func main() {
	flag.Parse()
	devMode := Mode == DEV_MODE

	if devMode {
		err := godotenv.Load("local.env")
		if err != nil {
			panic(err)
		}
	}
	var listenAddress = os.Getenv("LISTEN_ADDRESS")
	var successRedirect = os.Getenv("SUCCESS_REDIRECT")
	var errorRedirect = os.Getenv("ERROR_REDIRECT")

	var authServiceHost = os.Getenv("AUTH_SERVICE_HOST")
	var mailServerHost = os.Getenv("MAIL_SERVER_HOST")
	var mailServerPort = os.Getenv("MAIL_SERVER_PORT")
	var fromAddress = os.Getenv("MAIL_FROM_ADDRESS")

	sessions := cache.New(5*time.Minute, 10*time.Minute)

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
