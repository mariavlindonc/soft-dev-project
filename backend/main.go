package main

import (
	"log"

	db "backend/dao"
	"backend/domain"
)

func main() {
	dbase, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer dbase.Close()

	if err := dbase.Migrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
