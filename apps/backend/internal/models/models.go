package models

import "time"

type Order struct {
	orderUID          [16]byte
	trackNumber       string
	entry             string
	deliveryID        [16]byte
	paymentID         [16]byte
	locale            string
	internalSignature string
	customerID        string
	deliveryService   string
	shardkey          string
	smID              int
	dateCreated       time.Time
	oofShard          string
}

type Delivery struct {
	id      [16]byte
	name    string
	phone   string
	zip     string
	city    string
	address string
	region  string
	email   string
}

type Payment struct {
	id           [16]byte
	transaction  string
	requiestID   string
	currency     string
	provider     string
	amount       int
	paymentDT    int
	bank         string
	deliveryCost int
	goodsTotal   int
	customFee    int
}

type Item struct {
	id          [16]byte
	orderUID    [16]byte
	chrtID      int
	trackNumber string
	price       int
	rid         string
	name        string
	sale        int
	size        string
	totalPrice  int
	nmID        int
	brand       string
	status      int
}
