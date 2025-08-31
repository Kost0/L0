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

	// define signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// goroutine for server
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.StartHTTPServer(ctx, repo, orderCache)
	}()

	// goroutine for kafka
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.StartKafka(ctx, repo)
	}()

	wg.Wait()
	log.Println("All components stopped gracefully")
}
