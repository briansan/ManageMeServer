package main

import (
	"os"

	"github.com/briansan/ManageMeServer/api"
	"github.com/briansan/ManageMeServer/model/store"
	"github.com/briansan/ManageMeServer/www"
)

const (
	envMode = "MANAGEME_SERVER_MODE"
)

func main() {
	if mode := os.Getenv(envMode); mode == "api" {
		err := store.InitMongoSession()
		if err != nil {
			panic(err)
		}

		api.New().Start(":8888")
	} else {
		www.New().Start(":8889")
	}
}
