package JSONAPI

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mailing-list/dataBase"
	"net/http"
)

func setJSONHeader(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func fromJSON[T any](body io.Reader, target T) {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(body)
	json.Unmarshal(buffer.Bytes(), &target)
}

func returnJSON[T any](writer http.ResponseWriter, withData func() (T, error)) {
	setJSONHeader(writer)
	data, serverError := withData()

	if serverError != nil {
		writer.WriteHeader(500)
		serverErrorJSON, err := json.Marshal(&serverError)

		if err != nil {
			log.Print(err)

			return
		} else {
			writer.Write(serverErrorJSON)

			return
		}
	}

	dataJSON, err := json.Marshal(&data)

	if err != nil {
		log.Print(err)
		writer.WriteHeader(500)

		return
	} else {
		writer.Write(dataJSON)
	}
}

func returnError(writer http.ResponseWriter, err error, responseCode int) {
	returnJSON(writer, func() (interface{}, error) {
		errorMessage := struct {
			Err string
		}{
			Err: err.Error(),
		}
		writer.WriteHeader(responseCode)

		return errorMessage, nil
	})
}

func CreateEMail(eMailDataBase *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			return
		}

		eMail := dataBase.EMail{}
		fromJSON(request.Body, &eMail)

		if err := dataBase.CreateEMail(eMailDataBase, eMail.Value); err != nil {
			returnError(writer, err, 400)

			return
		} else {
			returnJSON(writer, func() (interface{}, error) {
				log.Printf("JSON-CreateEMail: %v\n", eMail.Value)

				return dataBase.GetEMail(eMailDataBase, eMail.Value)
			})
		}
	})
}

func GetEMail(eMailDataBase *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			return
		}

		eMail := dataBase.EMail{}
		fromJSON(request.Body, &eMail)

		returnJSON(writer, func() (interface{}, error) {
			log.Printf("JSON-GetEMail: %v\n", eMail.Value)

			return dataBase.GetEMail(eMailDataBase, eMail.Value)
		})
	})
}

func GetEMailBatch(eMailDataBase *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			return
		}

		parameters := dataBase.GetEMailBatchParameters{}
		fromJSON(request.Body, &parameters)

		if parameters.Count <= 0 || parameters.Page <= 0 {
			returnError(writer, errors.New("inValid count & page"), 400)

			return
		} else {
			returnJSON(writer, func() (interface{}, error) {
				log.Printf("JSON-GetEMailBatch: %v\n", parameters)

				return dataBase.GetEMailBatch(eMailDataBase, parameters)
			})
		}
	})
}

func UpDateEMail(eMailDataBase *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "PUT" {
			return
		}

		eMail := dataBase.EMail{}
		fromJSON(request.Body, &eMail)

		if err := dataBase.UpDateEMail(eMailDataBase, eMail); err != nil {
			returnError(writer, err, 400)

			return
		} else {
			returnJSON(writer, func() (interface{}, error) {
				log.Printf("JSON-UpDateEMail: %v\n", eMail.Value)

				return dataBase.GetEMail(eMailDataBase, eMail.Value)
			})
		}
	})
}

func DeleteEMail(eMailDataBase *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "DELETE" {
			return
		}

		eMail := dataBase.EMail{}
		fromJSON(request.Body, &eMail)

		if err := dataBase.DeleteEMail(eMailDataBase, eMail.Value); err != nil {
			returnError(writer, err, 400)

			return
		} else {
			returnJSON(writer, func() (interface{}, error) {
				log.Printf("JSON-DeleteEMail: %v\n", eMail.Value)

				return dataBase.GetEMail(eMailDataBase, eMail.Value)
			})
		}
	})
}

func Serve(eMailDataBase *sql.DB, address string) {
	http.Handle("/email/create", CreateEMail(eMailDataBase))
	http.Handle("/email/get", GetEMail(eMailDataBase))
	http.Handle("/email/get_batch", GetEMailBatch(eMailDataBase))
	http.Handle("/email/update", UpDateEMail(eMailDataBase))
	http.Handle("/email/delete", DeleteEMail(eMailDataBase))

	log.Printf("JSON-server listening on: %v\n", address)

	err := http.ListenAndServe(address, nil)

	if err != nil {
		log.Fatalf("JSON-server--error: %v", err)
	}
}