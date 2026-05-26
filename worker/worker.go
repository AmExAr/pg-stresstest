package worker

import (
	"context"
	"log"
	"math/rand"
	"pg-stresstest/db"
	"pg-stresstest/model"
	"sync"
	"time"
)

type RecordsDB struct {
	Records map[int]model.Record
	Mu      sync.Mutex
}

func Worker(connString string, rdb *RecordsDB, idx chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := db.ConnectDB(connString)
	if err != nil {
		log.Fatalf("Поток не установил соединение с БД: %s\n", err)
	}
	defer conn.Close(context.Background())

	for id := range idx {
		if conn == nil {
			rdb.Mu.Lock()
			rdb.Records[id] = model.Record{
				ID:     id,
				Sent:   false,
				Exists: false,
			}
			rdb.Mu.Unlock()
			conn, _ = db.RecreateConnection(conn, connString)
			continue
		}
		if err := conn.Ping(context.Background()); err != nil {
			rdb.Mu.Lock()
			rdb.Records[id] = model.Record{
				ID:     id,
				Sent:   false,
				Exists: false,
			}
			rdb.Mu.Unlock()
			conn, _ = db.RecreateConnection(conn, connString)
			continue
		}

		err := db.InsertID(conn, id)

		rdb.Mu.Lock()
		rdb.Records[id] = model.Record{
			ID:     id,
			Sent:   err == nil,
			Exists: false,
		}
		rdb.Mu.Unlock()

		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
}
