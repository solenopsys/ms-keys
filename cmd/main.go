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

const SECRET_NAME = "platform-ms-keys"
const DEV_MODE = "dev"

func init() {
	flag.StringVar(&Mode, "mode", "", "a string var")
}

func main() {
	flag.Parse()
	devMode := Mode == DEV_MODE

	if devMode {
		err := godotenv.Load("local-sec.env")
		if err != nil {
			panic(err)
		}
	}
	var dbPath = os.Getenv("DB_PATH")
	var listenAddress = os.Getenv("LISTEN_ADDRESS")
	var successRedirect = os.Getenv("SUCCESS_REDIRECT")
	var errorRedirect = os.Getenv("ERROR_REDIRECT")

	var authServiceHost = os.Getenv("AUTH_SERVICE_HOST")
	var mailServerHost = os.Getenv("MAIL_SERVER_HOST")
	var mailServerPort = os.Getenv("MAIL_SERVER_PORT")
	var fromAddress = os.Getenv("MAIL_FROM_ADDRESS")
	var smtpUser = os.Getenv("SMTP_USERNAME")
	var smtpPass = os.Getenv("SMTP_PASSWORD")

	sessions := cache.New(5*time.Minute, 10*time.Minute)

	th := internal.NewTransportHolder()
	if devMode {
		th.AddTransport("log", &internal.LogTransport{})
	}

	mt := internal.MailTransport{
		AuthHost: authServiceHost,
		From:     fromAddress,
		Host:     mailServerHost,
		Port:     mailServerPort,
		Password: string(smtpPass),
		Username: string(smtpUser),
	}
	th.AddTransport("email", &mt)

	db := internal.NewDb(dbPath)

	restServer := internal.RestServer{
		Sessions:         sessions,
		TransportService: th,
		SuccessUrl:       successRedirect,
		ErrorUrl:         errorRedirect,
		ListenAddress:    listenAddress,
		Db:               db,
	}

	println("start server on " + listenAddress)
	restServer.Run()

	defer db.Close()
}
