package gRPCAPI

import (
	"context"
	"database/sql"
	"log"
	"mailing-list/dataBase"
	protoBuf "mailing-list/protoBuf"
	"net"
	"time"

	"google.golang.org/grpc"
)

type EMailServer struct {
	protoBuf.UnimplementedMailingListServiceServer
	eMailDataBase *sql.DB
}

func protoBufEMailToDataBaseEMail(protoBufEMail *protoBuf.EMail) dataBase.EMail {
	confirmedAtTime := time.Unix(protoBufEMail.ConfirmedAt, 0)

	return dataBase.EMail{
		ID:          protoBufEMail.ID,
		Value:       protoBufEMail.Value,
		ConfirmedAt: &confirmedAtTime,
		OptOut:      protoBufEMail.OptOut,
	}
}

func dataBaseEMailToProtoBufEMail(dataBaseEMail *dataBase.EMail) protoBuf.EMail {
	return protoBuf.EMail{
		ID:          dataBaseEMail.ID,
		Value:       dataBaseEMail.Value,
		ConfirmedAt: dataBaseEMail.ConfirmedAt.Unix(),
		OptOut:      dataBaseEMail.OptOut,
	}
}

func eMailResponse(eMailDataBase *sql.DB, eMail string) (*protoBuf.EMailResponse, error) {
	dataBaseEMail, err := dataBase.GetEMail(eMailDataBase, eMail)

	if err != nil {
		return &protoBuf.EMailResponse{}, err
	}

	if dataBaseEMail == nil {
		return &protoBuf.EMailResponse{}, nil
	}

	protoBufEMail := dataBaseEMailToProtoBufEMail(dataBaseEMail)

	return &protoBuf.EMailResponse{EMail: &protoBufEMail}, nil
}

func (eMailServer *EMailServer) CreateEMail(context context.Context, request *protoBuf.CreateEMailReQuest) (*protoBuf.EMailResponse, error) {
	log.Printf("gRPC-CreateEMail: %v\n", request)

	err := dataBase.CreateEMail(eMailServer.eMailDataBase, request.Value)

	if err != nil {
		return &protoBuf.EMailResponse{}, err
	} else {
		return eMailResponse(eMailServer.eMailDataBase, request.Value)
	}
}

func (eMailServer *EMailServer) GetEMail(context context.Context, request *protoBuf.GetEMailReQuest) (*protoBuf.EMailResponse, error) {
	log.Printf("gRPC-GetEMail: %v\n", request)

	return eMailResponse(eMailServer.eMailDataBase, request.Value)
}

func (eMailServer *EMailServer) GetEMailBatch(context context.Context, request *protoBuf.GetEMailBatchReQuest) (*protoBuf.EMailBatchResponse, error) {
	log.Printf("gRPC-GetEMailBatch: %v\n", request)

	parameters := dataBase.GetEMailBatchParameters{
		Count: int(request.Count),
		Page:  int(request.Page),
	}

	dataBaseEMails, err := dataBase.GetEMailBatch(eMailServer.eMailDataBase, parameters)

	if err != nil {
		return &protoBuf.EMailBatchResponse{}, err
	}

	protoBufEMails := make([]*protoBuf.EMail, 0, len(dataBaseEMails))

	for i := 0; i < len(dataBaseEMails); i++ {
		protoBufEMail := dataBaseEMailToProtoBufEMail(&dataBaseEMails[i])
		protoBufEMails = append(protoBufEMails, &protoBufEMail)
	}

	return &protoBuf.EMailBatchResponse{EMails: protoBufEMails}, nil
}

func (eMailServer *EMailServer) UpDateEMail(context context.Context, request *protoBuf.UpDateEMailReQuest) (*protoBuf.EMailResponse, error) {
	log.Printf("gRPC-UpDateEMail: %v\n", request)

	dataBaseEMail := protoBufEMailToDataBaseEMail(request.EMail)

	err := dataBase.UpDateEMail(eMailServer.eMailDataBase, dataBaseEMail)

	if err != nil {
		return &protoBuf.EMailResponse{}, err
	} else {
		return eMailResponse(eMailServer.eMailDataBase, dataBaseEMail.Value)
	}
}

func (eMailServer *EMailServer) DeleteEMail(context context.Context, request *protoBuf.DeleteEMailReQuest) (*protoBuf.EMailResponse, error) {
	log.Printf("gRPC-DeleteEMail: %v\n", request)

	err := dataBase.DeleteEMail(eMailServer.eMailDataBase, request.Value)

	if err != nil {
		return &protoBuf.EMailResponse{}, err
	} else {
		return eMailResponse(eMailServer.eMailDataBase, request.Value)
	}
}

func Serve(eMailDataBase *sql.DB, address string) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("gRPC-server--error: %v\n", address)
	}

	gRPCServer := grpc.NewServer()
	eMailServer := EMailServer{eMailDataBase: eMailDataBase}
	protoBuf.RegisterMailingListServiceServer(gRPCServer, &eMailServer)
	log.Printf("gRPC-server listening on: %v\n", address)

	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("gRPC-server--error: %v\n", err)
	}
}