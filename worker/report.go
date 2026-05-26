package worker

import (
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
	defer conn.Close(nil)

	records, err := db.FetchTable(conn)
	if err != nil {
		return fmt.Errorf("Не удалось получить данные из БД: %s\n", err)
	}

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
	defer writer.Flush()

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

	return nil
}
