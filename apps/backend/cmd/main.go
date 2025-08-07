package main

import (
	"github.com/Kost0/L0/apps/backend/internal/repository"
	"log"
)

func main() {
	db, err := repository.CreateDataBase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = repository.RunMigrations(db, "orders_l0")
	if err != nil {
		log.Fatal(err)
	}
}
