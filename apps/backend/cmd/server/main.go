package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Kost0/L0/internal/cache"
	"github.com/Kost0/L0/internal/http"
	"github.com/Kost0/L0/internal/kafka"
	"github.com/Kost0/L0/internal/repository"
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

	orderCache := cache.NewOrderCache(48 * time.Hour)

	err = orderCache.WarmUpCache(db, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.StartHTTPServer(ctx, db)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.StartKafka(ctx, db)
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
