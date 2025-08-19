package cache

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Kost0/L0/internal/repository"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(ctx context.Context) (*sql.DB, func(), error) {
	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	tearDown := func() {
		if err := container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate test container: %v", err)
		}
	}

	port, err := container.MappedPort(ctx, nat.Port("5432"))
	if err != nil {
		tearDown()
		return nil, nil, err
	}

	connStr := fmt.Sprintf("host=localhost port=%s user=test password=test dbname=test sslmode=disable", port.Port())

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		tearDown()
		return nil, nil, err
	}

	if err = db.Ping(); err != nil {
		tearDown()
		return nil, nil, err
	}

	return db, tearDown, nil
}

func createShema(db *sql.DB) error {
	schema := `
CREATE TABLE delivery (
    id VARCHAR(50),
    name VARCHAR(255),
    phone VARCHAR(50),
    zip VARCHAR(50),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE payment (
    id VARCHAR(50),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(3) ,
    provider VARCHAR(255),
    amount INT,
    payment_dt BIGINT,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(50),
    delivery_id VARCHAR(50),
    payment_id VARCHAR(50),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255),
    sm_id INT,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(255)
);

CREATE TABLE items (
    id VARCHAR(50),
    order_id VARCHAR(50),
    chrt_id INT,
    track_number VARCHAR(100),
    price INT,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(50),
    total_price INT,
    nm_id INT,
    brand VARCHAR(100),
    status INT
);
`

	_, err := db.Exec(schema)
	return err
}

func TestGetRecentOrders_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	db, tearDown, err := setupTestDB(ctx)
	assert.NoError(t, err)
	defer tearDown()
	defer db.Close()

	err = createShema(db)
	assert.NoError(t, err)

	_, _ = db.Exec("DELETE FROM orders")

	now := time.Now()

	_, _ = db.Exec("INSERT INTO orders (order_uid, date_created, delivery_id, payment_id) VALUES ($1, $2, $3, $4)", "order-new-1", now, "1", "2")
	_, _ = db.Exec("INSERT INTO orders (order_uid, date_created, delivery_id, payment_id) VALUES ($1, $2, $3, $4)", "order-new-2", now.Add(-48*time.Hour), "1", "2")
	_, _ = db.Exec("INSERT INTO orders (order_uid, date_created, delivery_id, payment_id) VALUES ($1, $2, $3, $4)", "order-old-1", now.Add(-240*time.Hour), "1", "2")

	_, _ = db.Exec("INSERT INTO delivery (id) VALUES ($1)", "1")
	_, _ = db.Exec("INSERT INTO payment (id) VALUES ($1)", "2")
	_, _ = db.Exec("INSERT INTO items (order_id) VALUES ($1)", "order-new-1")
	_, _ = db.Exec("INSERT INTO items (order_id) VALUES ($1)", "order-new-2")
	_, _ = db.Exec("INSERT INTO items (order_id) VALUES ($1)", "order-old-1")

	repo := repository.NewOrderRepository(db)

	data, err := getRecentOrders(db, repo, ctx)

	assert.NoError(t, err)
	assert.Len(t, data, 2)

	orders := make(map[string]bool)
	for _, o := range data {
		orders[o.Order.OrderUID] = true
	}
	assert.True(t, orders["order-new-1"])
	assert.True(t, orders["order-new-2"])
	assert.False(t, orders["order-old-1"])
}
