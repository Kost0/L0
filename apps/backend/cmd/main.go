package main

import (
	"context"
	"github.com/Kost0/L0/internal/repository"
	"github.com/Kost0/L0/internal/server"
	"github.com/Kost0/L0/internal/workWithKafka"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting server")
	defer db.Close()

	err = repository.RunMigrations(db, "orders_l0")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrations complete")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartHTTPServer(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		workWithKafka.StartKafka(ctx)
	}()

	<-done
	log.Println("Shutting down")

	cancel()

	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		log.Println("All components stopped gracefully")
	case <-time.After(time.Second * 10):
		log.Println("All components failed to stop gracefully")
	}
}
