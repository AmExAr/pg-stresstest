package worker

import (
	"log"
	"math/rand"
	"pg-stresstest/db"
	"pg-stresstest/model"
	"sync"
	"time"
)

type RecordsDB struct {
	Records map[int]model.Record
	Mu      *sync.Mutex
}

func Worker(connString string, rdb *RecordsDB, idx chan int) {
	conn, err := db.ConnectDB(connString)
	if err != nil {
		log.Fatalf("Поток не установил соединение с БД: %s\n", err)
	}
	defer conn.Close(nil)

	for {
		id := <-idx

		rdb.Mu.Lock()
		err := db.InsertID(conn, id)

		if err != nil {
			rdb.Records[id] = model.Record{
				ID:     id,
				Sent:   false,
				Exists: false,
			}
		} else {
			rdb.Records[id] = model.Record{
				ID:     id,
				Sent:   true,
				Exists: false,
			}
		}
		rdb.Mu.Unlock()

		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		conn, _ := db.RecreateConnection(conn, connString)
	}
}
