package main

import (
	"github.com/khrees2412/evolvecredit/database"
	"log"
	"os"
)

func Server() {
	dbUrl := os.Getenv("DATABASE_URL")
	dbConnection, err := database.ConnectDB(dbUrl)

	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	err = database.MigrateAll(dbConnection)

	if err != nil {
		log.Fatalf("migration error: %v", err)
	}
}
