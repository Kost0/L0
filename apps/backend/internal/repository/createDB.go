package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func CreateDataBase() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=postgres password=%s host=localhost sslmode=disable", os.Getenv("DB_PASSWORD"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE orders_l0")
	if err != nil {
		fmt.Printf("Database exists or error: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("Database created successfully\n")
	}

	_, err = db.Exec(fmt.Sprintf("CREATE USER orders_user WITH PASSWORD '%s'", os.Getenv("NEW_USER_PASSWORD")))
	if err != nil {
		fmt.Printf("User exists or error: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("User created successfully\n")
	}

	_, err = db.Exec("GRANT ALL PRIVILEGES ON DATABASE orders_l0 To orders_user")
	if err != nil {
		fmt.Printf("User exists or error: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("User granted successfully\n")
	}

	adminConnStr := fmt.Sprintf("user=postgres password=%s host=localhost dbname=orders_l0 sslmode=disable", os.Getenv("DB_PASSWORD"))
	db, err = sql.Open("postgres", adminConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to new database as admin: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("GRANT ALL PRIVILEGES ON SCHEMA public TO orders_user")
	if err != nil {
		fmt.Printf("User exists or error: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("User granted successfully\n")
	}

	_, err = db.Exec("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO orders_user")
	if err != nil {
		fmt.Printf("User exists or error: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("User granted successfully\n")
	}

	_, err = db.Exec("GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO orders_user")
	if err != nil {
		fmt.Printf("User exists or error: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("User granted successfully\n")
	}

	db, err = sql.Open("postgres", fmt.Sprintf("user=orders_user password=%s host=localhost sslmode=disable dbname=orders_l0", os.Getenv("NEW_USER_PASSWORD")))
	if err != nil {
		fmt.Printf("Database exists or error: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("Database opened successfully\n")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
