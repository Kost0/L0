// Package models provides structs of all data
package models

import (
	"time"
)

// Order provides main information about the order
// @Description Main order information
type Order struct {
	OrderUID          string  `fake:"{uuid}"`
	TrackNumber       *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Entry             *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	DeliveryID        *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Locale            *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	InternalSignature *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	CustomerID        *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	DeliveryService   *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Shardkey          *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	SmID              *int    `fake:"{uint8}"`
	DateCreated       *time.Time
	OofShard          *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
}

// Delivery provides information about the delivery
// @Description delivery information
type Delivery struct {
	ID      *string `fake:"{uuid}"`
	Name    *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Phone   *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Zip     *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	City    *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Address *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Region  *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Email   *string `fake:"{email}"`
}

// Payment provides information about the payment
// @Description payment information
type Payment struct {
	Transaction  *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	RequestID    *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Currency     *string `fake:"{regex:[a-zA-Z0-9]{1,3}}"`
	Provider     *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Amount       *int    `fake:"{uint8}"`
	PaymentDT    *int    `fake:"{uint8}"`
	Bank         *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	DeliveryCost *int    `fake:"{uint8}"`
	GoodsTotal   *int    `fake:"{uint8}"`
	CustomFee    *int    `fake:"{uint8}"`
}

// Item provides information about the product
// @Description Product information
type Item struct {
	ChrtID      *int    `fake:"{uint8}"`
	TrackNumber *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Price       *int    `fake:"{uint8}"`
	Rid         *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Name        *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Sale        *int    `fake:"{uint8}"`
	Size        *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	TotalPrice  *int    `fake:"{uint8}"`
	NmID        *int    `fake:"{uint8}"`
	Brand       *string `fake:"{regex:[a-zA-Z0-9]{1,10}}"`
	Status      *int    `fake:"{uint8}"`
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
