package models

import (
	"database/sql"
)

// Order provides main information about the order
// @Description Main order information
type Order struct {
	OrderUID          string       `json:"orderUID"`
	TrackNumber       *string      `json:"trackNumber"`
	Entry             *string      `json:"entry"`
	DeliveryID        *string      `json:"deliveryID"`
	PaymentID         *string      `json:"paymentID"`
	Locale            *string      `json:"locale"`
	InternalSignature *string      `json:"internalSignature"`
	CustomerID        *string      `json:"customerID"`
	DeliveryService   *string      `json:"deliveryService"`
	Shardkey          *string      `json:"shardKey"`
	SmID              *int         `json:"smID"`
	DateCreated       sql.NullTime `json:"dateCreated"`
	OofShard          *string      `json:"oofShard"`
}

// Delivery provides information about the delivery
// @Description delivery information
type Delivery struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Zip     *string `json:"zip"`
	City    *string `json:"city"`
	Address *string `json:"address"`
	Region  *string `json:"region"`
	Email   *string `json:"email"`
}

// Payment provides information about the payment
// @Description payment information
type Payment struct {
	ID           *string `json:"id"`
	Transaction  *string `json:"transaction"`
	RequestID    *string `json:"requestID"`
	Currency     *string `json:"currency"`
	Provider     *string `json:"provider"`
	Amount       *int    `json:"amount"`
	PaymentDT    *int    `json:"paymentDT"`
	Bank         *string `json:"bank"`
	DeliveryCost *int    `json:"deliveryCost"`
	GoodsTotal   *int    `json:"goodsTotal"`
	CustomFee    *int    `json:"customFee"`
}

// Item provides information about the product
// @Description Product information
type Item struct {
	ID          *string `json:"id"`
	OrderUID    *string `json:"orderUID"`
	ChrtID      *int    `json:"chrtID"`
	TrackNumber *string `json:"trackNumber"`
	Price       *int    `json:"price"`
	Rid         *string `json:"rid"`
	Name        *string `json:"name"`
	Sale        *int    `json:"sale"`
	Size        *string `json:"size"`
	TotalPrice  *int    `json:"totalPrice"`
	NmID        *int    `json:"nmID"`
	Brand       *string `json:"brand"`
	Status      *int    `json:"status"`
}

// CombinedData presents all the information about the order
// @Description Information about the order and nested structures
type CombinedData struct {
	// Main order information
	Order Order
	// payment information
	Payment Payment
	// delivery information
	Delivery Delivery
	// items information
	Items []Item
}
