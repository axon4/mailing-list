package main

import (
	"context"
	"log"
	protoBuf "mailing-list/protoBuf"
	"time"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func logResponse(response *protoBuf.EMailResponse, err error) {
	if err != nil {
		log.Fatalf("	error: %v", err)
	}

	if response.EMail == nil {
		log.Printf("	eMail not found")
	} else {
		log.Printf("	response: %v", response.EMail)
	}
}

func createEMail(client protoBuf.MailingListServiceClient, value string) *protoBuf.EMail {
	log.Println("client-CreateEMail")

	eMailConText, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	response, err := client.CreateEMail(eMailConText, &protoBuf.CreateEMailReQuest{Value: value})
	logResponse(response, err)

	return response.EMail
}

func getEMail(client protoBuf.MailingListServiceClient, value string) *protoBuf.EMail {
	log.Println("client-GetEMail")

	eMailConText, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	response, err := client.GetEMail(eMailConText, &protoBuf.GetEMailReQuest{Value: value})
	logResponse(response, err)

	return response.EMail
}

func getEMailBatch(client protoBuf.MailingListServiceClient, count int, page int) {
	log.Println("client-GetEMailBatch")

	eMailConText, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	response, err := client.GetEMailBatch(eMailConText, &protoBuf.GetEMailBatchReQuest{Count: int32(count), Page: int32(page)})

	if err != nil {
		log.Fatalf("	error: %v", err)
	}

	log.Println("	response:")

	for i := 0; i < len(response.EMails); i++ {
		log.Printf("		item %v of %v | %s", i+1, len(response.EMails), response.EMails[i])
	}
}

func upDateEMail(client protoBuf.MailingListServiceClient, eMail *protoBuf.EMail) *protoBuf.EMail {
	log.Println("client-UpDateEMail")

	eMailConText, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	reQuest := protoBuf.UpDateEMailReQuest{EMail: eMail}
	response, err := client.UpDateEMail(eMailConText, &reQuest)
	logResponse(response, err)

	return response.EMail
}

func deleteEMail(client protoBuf.MailingListServiceClient, value string) *protoBuf.EMail {
	log.Println("client-DeleteEMail")

	eMailConText, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	response, err := client.DeleteEMail(eMailConText, &protoBuf.DeleteEMailReQuest{Value: value})
	logResponse(response, err)

	return response.EMail
}

var arguments struct {
	gRPCAddress string `arg:"env:MAILING_LIST___GRPC_SERVER__ADDRESS"`
}

func main() {
	arg.MustParse(&arguments)

	if arguments.gRPCAddress == "" {
		arguments.gRPCAddress = ":3002"
	}

	connection, err := grpc.Dial(arguments.gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("gRPC-client--error: %v", err)
	}

	defer connection.Close()
	client := protoBuf.NewMailingListServiceClient(connection)
	const TestEMail = "test@test.test";

	eMail := createEMail(client, TestEMail)
	eMail.ConfirmedAt = 77777
	upDateEMail(client, eMail)
	deleteEMail(client, eMail.Value)
	getEMail(client, TestEMail)
	getEMailBatch(client, 3, 1)
	getEMailBatch(client, 3, 2)
	getEMailBatch(client, 3, 3)
}