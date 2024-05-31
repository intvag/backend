package main

import (
	"os"

	"github.com/gofrs/uuid/v5"
	decisions "github.com/intvag/decision-engine/service"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	decisionsEngineAddr = envOrDefault("BACKEND_DECISIONS_ENGINE", "localhost:8888")
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Person{}, &Item{}, &Policy{}, &PolicyItem{}, &Quote{}, &QuoteItem{})

	db.Create(&Item{
		ID:              uuid.Must(uuid.NewV4()),
		Category:        "Fridge",
		Manufacturer:    "Samsung",
		ModelName:       "CoolFridge",
		AverageLifetime: 20,
		Lastability:     0.33,
		Repairability:   0.66,
	})

	conn, err := grpc.Dial(decisionsEngineAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	qc := decisions.NewQuotesClient(conn)

	s, err := New(db, qc)
	if err != nil {
		panic(err)
	}

	panic(s.r.Run(":8989"))
}

func envOrDefault(k, d string) string {
	v, ok := os.LookupEnv(k)
	if ok {
		return v
	}

	return d
}
