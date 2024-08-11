package main

import (
	"database/sql"
	"log"
	JSONAPI "mailing-list/API/JSON"
	"mailing-list/dataBase"
	"sync"

	"github.com/alexflint/go-arg"
)

var arguments struct {
	dataBasePath      string `arg:"env:MAILING_LIST__DATABASE"`
	JSONServerAddress string `arg:"env:MAILING_LIST___JSON_SERVER__ADDRESS"`
}

func main() {
	arg.MustParse(&arguments)

	if arguments.dataBasePath == "" {
		arguments.dataBasePath = "dataBase.db"
	}

	if arguments.JSONServerAddress == "" {
		arguments.JSONServerAddress = ":3001"
	}

	eMailDataBase, err := sql.Open("sqlite3", arguments.dataBasePath)

	if err != nil {
		log.Fatal(err)
	}

	defer eMailDataBase.Close()

	dataBase.CreateTable(eMailDataBase)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("starting JSON-server\n")
		JSONAPI.Serve(eMailDataBase, arguments.JSONServerAddress)
		wg.Done()
	}()

	wg.Wait()
}