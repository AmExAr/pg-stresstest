package main

import (
	"context"
	"log"
	"os"
	"pg-stresstest/db"
	"pg-stresstest/model"
	"pg-stresstest/worker"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func main() {
	connString := os.Getenv("DB_CONN_STRING")
	THREADS := 20
	ITERATIONS := 10000

	conn, err := db.ConnectDB(connString)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %s\n", err)
	}
	defer conn.Close(context.Background())

	err = db.CreateTable(conn)
	conn.Close(context.Background())
	if err != nil {
		log.Fatalf("Ошибка создания/проверки таблицы в БД: %s\n", err)
	}

	rdb := &worker.RecordsDB{
		Records: make(map[int]model.Record),
	}
	idChan := make(chan int, 1)
	var wg sync.WaitGroup

	for range THREADS {
		wg.Add(1)
		go worker.Worker(connString, rdb, idChan, &wg)
	}

	bar := progressbar.Default(int64(ITERATIONS))
	for i := range ITERATIONS {
		idChan <- i
		bar.Add(1)
	}

	close(idChan)
	wg.Wait()

	err = worker.GenerateCSV(rdb, connString)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("А на этом всё. Отчёт о работе: report.csv")
}
