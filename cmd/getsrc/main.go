package main

import (
	"log"
	"net/http"

	"sheff.online/getsrc/internal/getsrc"
)

func main() {
	log.Println("Load config")
	config, err := getsrc.NewConfig("./getsrc.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Init HTTP")
	_, err = getsrc.NewHTTP(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Start HTTP server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
