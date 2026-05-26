package worker

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"pg-stresstest/db"
	"pg-stresstest/model"
	"strconv"
)

func GenerateCSV(rdb *RecordsDB, connString string) error {
	conn, err := db.ConnectDB(connString)
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к БД: %s\n", err)
	}
	defer conn.Close(context.Background())

	records, err := db.FetchTable(conn)
	if err != nil {
		return fmt.Errorf("Не удалось получить данные из БД: %s\n", err)
	}

	rdb.Mu.Lock()
	for _, r := range records {
		old := rdb.Records[r]
		rdb.Records[r] = model.Record{
			ID:     r,
			Sent:   old.Sent,
			Exists: true,
		}
	}

	file, err := os.Create("report.csv")
	if err != nil {
		return fmt.Errorf("Не удалось создать файл report.csv: %s\n", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	_ = writer.Write([]string{
		"id",
		"sent",
		"exists",
	})

	for _, r := range rdb.Records {
		row := []string{
			strconv.Itoa(r.ID),
			strconv.FormatBool(r.Sent),
			strconv.FormatBool(r.Exists),
		}

		_ = writer.Write(row)
	}

	rdb.Mu.Unlock()
	writer.Flush()
	return nil
}
