package main

import (
	"log"
	"os"

	server "github.com/Agilen/Server"
	"github.com/Agilen/Server/config"
)

func main() {
	os.Remove("DB.db")
	os.Create("DB.db")
	config := config.NewConfig("./config.conf")
	err := server.Start(config)
	if err != nil {
		log.Fatal(err)
	}
}
