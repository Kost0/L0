package models

import "time"

type Order struct {
	orderUID          [16]byte  `json:"orderUID"`
	trackNumber       string    `json:"trackNumber"`
	entry             string    `json:"entry"`
	deliveryID        [16]byte  `json:"deliveryID"`
	paymentID         [16]byte  `json:"paymentID"`
	locale            string    `json:"locale"`
	internalSignature string    `json:"internalSignature"`
	customerID        string    `json:"customerID"`
	deliveryService   string    `json:"deliveryService"`
	shardkey          string    `json:"shardKey"`
	smID              int       `json:"smID"`
	dateCreated       time.Time `json:"dateCreated"`
	oofShard          string    `json:"oofShard"`
}

type Delivery struct {
	id      [16]byte `json:"id"`
	name    string   `json:"name"`
	phone   string   `json:"phone"`
	zip     string   `json:"zip"`
	city    string   `json:"city"`
	address string   `json:"address"`
	region  string   `json:"region"`
	email   string   `json:"email"`
}

type Payment struct {
	id           [16]byte `json:"id"`
	transaction  string   `json:"transaction"`
	requestID    string   `json:"requestID"`
	currency     string   `json:"currency"`
	provider     string   `json:"provider"`
	amount       int      `json:"amount"`
	paymentDT    int      `json:"paymentDT"`
	bank         string   `json:"bank"`
	deliveryCost int      `json:"deliveryCost"`
	goodsTotal   int      `json:"goodsTotal"`
	customFee    int      `json:"customFee"`
}

type Item struct {
	id          [16]byte `json:"id"`
	orderUID    [16]byte `json:"orderUID"`
	chrtID      int      `json:"chrtID"`
	trackNumber string   `json:"trackNumber"`
	price       int      `json:"price"`
	rid         string   `json:"rid"`
	name        string   `json:"name"`
	sale        int      `json:"sale"`
	size        string   `json:"size"`
	totalPrice  int      `json:"totalPrice"`
	nmID        int      `json:"nmID"`
	brand       string   `json:"brand""`
	status      int      `json:"status"`
}
