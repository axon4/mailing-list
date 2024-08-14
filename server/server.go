package main

import (
	"database/sql"
	"log"
	JSONAPI "mailing-list/API/JSON"
	gRPCAPI "mailing-list/API/gRPC"
	"mailing-list/dataBase"
	"sync"

	"github.com/alexflint/go-arg"
)

var arguments struct {
	dataBase   string `arg:"env:MAILING_LIST__DATABASE"`
	JSONServer string `arg:"env:MAILING_LIST__JSON_SERVER"`
	gRPCServer string `arg:"env:MAILING_LIST__GRPC_SERVER"`
}

func main() {
	arg.MustParse(&arguments)

	if arguments.dataBase == "" {
		arguments.dataBase = "dataBase.db"
	}

	if arguments.JSONServer == "" {
		arguments.JSONServer = ":3001"
	}

	if arguments.gRPCServer == "" {
		arguments.gRPCServer = ":3002"
	}

	eMailDataBase, err := sql.Open("sqlite3", arguments.dataBase)

	if err != nil {
		log.Fatal(err)
	}

	defer eMailDataBase.Close()
	dataBase.CreateTable(eMailDataBase)

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		log.Printf("starting JSON-server\n")
		JSONAPI.Serve(eMailDataBase, arguments.JSONServer)
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		log.Printf("starting gRPC-server\n")
		gRPCAPI.Serve(eMailDataBase, arguments.gRPCServer)
		waitGroup.Done()
	}()

	waitGroup.Wait()
}