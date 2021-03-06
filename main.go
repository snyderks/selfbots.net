package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/snyderks/selfbots.net/configRead"
	"github.com/snyderks/selfbots.net/handlers"
)

func main() {
	config, err := configRead.ReadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	handlers.SetUpAPICalls()
	handlers.SetUpBasicHandlers()
	svr := http.Server{
		Addr:           config.HTTPPort,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 8175, // if it's good enough for Apache, it's good enough for me
	}
	fmt.Println("Serving", config.Hostname, "on", config.HTTPPort)
	svr.ListenAndServe()
}
