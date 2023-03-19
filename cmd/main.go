package main

import (
	"flag"
	_ "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"ms-keys/internal"
	"ms-keys/register-transport"
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

	th := register_transport.NewTransportHolder()
	if (devMode) {
		th.AddTransport("log", &register_transport.LogTransport{});
	}
	mt := register_transport.MailTransport{
		AuthHost: authServiceHost,
		From:     fromAddress,
		Host:     mailServerHost,
		Port:     mailServerPort,
	}
	th.AddTransport("mail", &mt)

	restServer := internal.RestServer{
		Sessions:         sessions,
		TransportService: th,
		SuccessUrl:       successRedirect,
		ErrorUrl:         errorRedirect,
		ListenAddress:    listenAddress,
	}

	println("start server on " + listenAddress)
	restServer.Run()
}
