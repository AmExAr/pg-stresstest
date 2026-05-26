package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(connString string) (*pgx.Conn, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateTable(conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS stress_test (
		id BIGINT PRIMARY KEY
	);
	`

	_, err := conn.Exec(context.Background(), query)
	return err
}

func InsertID(conn *pgx.Conn, id int) error {
	query := `INSERT INTO stress_test(id) VALUES($1)`

	_, err := conn.Exec(context.Background(), query, id)
	return err
}

func RecreateConnection(old *pgx.Conn, connString string) (*pgx.Conn, error) {
	if old != nil {
		_ = old.Close(context.Background())
	}

	log.Println("Переподключение к бд...")
	return ConnectDB(connString)
}

func FetchTable(conn *pgx.Conn) ([]int, error) {
	query := `SELECT id FROM stress_test`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []int{}

	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		records = append(records, id)
	}

	return records, nil
}
