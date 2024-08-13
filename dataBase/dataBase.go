package dataBase

import (
	"database/sql"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

type EMail struct {
	ID          int64
	Value       string
	ConfirmedAt *time.Time
	OptOut      bool
}

func CreateTable(dataBase *sql.DB) {
	_, err := dataBase.Exec(`
		CREATE TABLE eMails (
			ID            INTEGER PRIMARY KEY,
			eMail         TEXT UNIQUE,
			confirmed_at  INTEGER,
			opt_out       INTEGER
		);
	`)

	if err != nil {
		if SQLError, OK := err.(sqlite3.Error); OK {
			if SQLError.Code != 1 { // error-code 1: table already exists
				log.Fatal(SQLError)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func getEMailFromRow(row *sql.Rows) (*EMail, error) {
	var ID int64
	var value string
	var confirmedAt int64
	var optOut bool

	err := row.Scan(&ID, &value, &confirmedAt, &optOut)

	if err != nil {
		log.Println(err)

		return nil, err
	} else {
		confirmedAtTime := time.Unix(confirmedAt, 0)

		return &EMail{ID: ID, Value: value, ConfirmedAt: &confirmedAtTime, OptOut: optOut}, nil
	}
}

func CreateEMail(dataBase *sql.DB, eMail string) error {
	_, err := dataBase.Exec(`INSERT INTO eMails (email, confirmed_at, opt_out) VALUES(?, 0, false)`, eMail)

	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func GetEMail(dataBase *sql.DB, eMail string) (*EMail, error) {
	rows, err := dataBase.Query(`SELECT ID, eMail, confirmed_at, opt_out FROM eMails WHERE eMail = ?`, eMail)

	if err != nil {
		log.Println(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		return getEMailFromRow(rows)
	}

	return nil, nil
}

type GetEMailBatchParameters struct {
	Count int
	Page  int
}

func GetEMailBatch(dataBase *sql.DB, parameters GetEMailBatchParameters) ([]EMail, error) {
	var empty []EMail

	rows, err := dataBase.Query(`SELECT ID, eMail, confirmed_at, opt_out FROM eMails WHERE opt_out = false ORDER BY ID ASC LIMIT ? OFFSET ?`, parameters.Count, (parameters.Page-1)*parameters.Count)

	if err != nil {
		log.Println(err)

		return empty, err
	}

	defer rows.Close()

	eMails := make([]EMail, 0, parameters.Count)

	for rows.Next() {
		email, err := getEMailFromRow(rows)

		if err != nil {
			return nil, err
		}

		eMails = append(eMails, *email)
	}

	return eMails, nil
}

func UpDateEMail(dataBase *sql.DB, eMail EMail) error {
	confirmedAtTime := eMail.ConfirmedAt.Unix()

	_, err := dataBase.Exec(`INSERT INTO eMails (eMail, confirmed_at, opt_out) VALUES (?, ?, ?) ON CONFLICT (eMail) DO UPDATE SET confirmed_at=?, opt_out=?`, eMail.Value, confirmedAtTime, eMail.OptOut, confirmedAtTime, eMail.OptOut)

	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}

func DeleteEMail(dataBase *sql.DB, eMail string) error {
	_, err := dataBase.Exec(`UPDATE eMails SET opt_out=true WHERE eMail=?`, eMail)

	if err != nil {
		log.Println(err)

		return err
	} else {
		return nil
	}
}