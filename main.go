package main

import (
	"log"
	"pg-stresstest/db"
	"pg-stresstest/model"
	"pg-stresstest/worker"
)

func main() {
	connString := "postgres://postgres:1234567890@localhost:5432/postgres"
	THREADS := 10

	conn, err := db.ConnectDB(connString)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %s\n", err)
	}
	defer conn.Close(nil)

	err = db.CreateTable(conn)
	if err != nil {
		log.Fatalf("Ошибка создания/проверки таблицы в БД: %s\n", err)
	}

	rdb := worker.RecordsDB{
		Records: make(map[int]model.Record),
	}
	idChan := make(chan int)

	for i := 0; i < THREADS; i++ {
		go worker.Worker(connString, &rdb, idChan)
	}

	for i := 0; i < 1000000; i++ {
		idChan <- i
	}

	err = worker.GenerateCSV(&rdb, connString)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("А на этом всё. Отчёт о работе: report.csv")
}
