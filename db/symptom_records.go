package db

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"log"
)

type SymptomsRecord struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Symptoms string `json:"symptoms"`
	Disease  string `json:"disease"`
	UserId   int    `json:"user_id"`
	Age      int    `json:"age"`
	Rating   int    `json:"rating"`
}

func GetAllRecords(connection string) ([]SymptomsRecord, error) {
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(conn, context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM Records")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []SymptomsRecord
	for rows.Next() {
		record := SymptomsRecord{}
		err = rows.Scan(&record.ID, &record.Date, &record.Symptoms, &record.Disease, &record.UserId, &record.Age, &record.Rating)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func GetRecordsByUserId(userId int) ([]SymptomsRecord, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close(context.Background())

	row, err := conn.Query(context.Background(), "SELECT record_id, date, symptoms, disease, user_id, age, rating FROM Records WHERE user_id=$1", userId)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var records []SymptomsRecord
	for row.Next() {
		record := SymptomsRecord{}
		err = row.Scan(&record.ID, &record.Date, &record.Symptoms, &record.Disease, &record.UserId, &record.Age, &record.Rating)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	err = row.Err()
	if err != nil {
		return nil, err
	}
	return records, nil
}
func InsertRecord(record SymptomsRecord) error {
	// Open a database connection
	db, err := sql.Open("postgres", "postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi")
	if err != nil {
		return err
	}
	defer db.Close()

	var id int
	err = db.QueryRow("SELECT nextval('records_id_seq')").Scan(&id)
	if err != nil {
		return err
	}

	// Execute the SQL INSERT statement
	stmt, err := db.Prepare("INSERT INTO Records(record_id, date, symptoms, disease, user_id, age, rating) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, record.Date, record.Symptoms, record.Disease, record.UserId, record.Age, record.Rating)
	if err != nil {
		return err
	}

	return nil
}
