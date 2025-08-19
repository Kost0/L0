// @title L0 API
// @version 1.0
// @description API for getting order information
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/Kost0/L0/docs"
	"github.com/Kost0/L0/internal/cache"
	"github.com/Kost0/L0/internal/http"
	"github.com/Kost0/L0/internal/kafka"
	"github.com/Kost0/L0/internal/repository"
)

func main() {
	//connecting to database
	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting server")
	defer db.Close()

	// start migrations
	err = repository.RunMigrations(db, "orders_l0")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrations complete")

	//create object to work with database
	repo := repository.NewOrderRepository(db)

	// create object to work with cache
	orderCache := cache.NewOrderCache(48 * time.Hour)

	// fills the cache with data from database
	err = orderCache.WarmUpCache(db, repo, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// make channel to work with some goroutines
	done := make(chan os.Signal, 1)
	// define signals for graceful shutdown
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// goroutine for server
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.StartHTTPServer(ctx, repo, orderCache)
	}()

	ctx2, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	// goroutine for cmd
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.StartKafka(ctx2, repo)
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
