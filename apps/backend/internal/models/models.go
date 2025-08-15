package models

import (
	"time"
)

// Order provides main information about the order
// @Description Main order information
type Order struct {
	OrderUID          string     `json:"orderUID" validate:"required"`
	TrackNumber       *string    `json:"trackNumber" validate:"required"`
	Entry             *string    `json:"entry" validate:"required"`
	DeliveryID        *string    `json:"deliveryID" validate:"required"`
	PaymentID         *string    `json:"paymentID" validate:"required"`
	Locale            *string    `json:"locale" validate:"required"`
	InternalSignature *string    `json:"internalSignature"`
	CustomerID        *string    `json:"customerID" validate:"required"`
	DeliveryService   *string    `json:"deliveryService" validate:"required"`
	Shardkey          *string    `json:"shardKey" validate:"required"`
	SmID              *int       `json:"smID" validate:"required"`
	DateCreated       *time.Time `json:"dateCreated" validate:"required"`
	OofShard          *string    `json:"oofShard" validate:"required"`
}

// Delivery provides information about the delivery
// @Description delivery information
type Delivery struct {
	ID      *string `json:"id" validate:"required"`
	Name    *string `json:"name" validate:"required"`
	Phone   *string `json:"phone" validate:"required"`
	Zip     *string `json:"zip" validate:"required"`
	City    *string `json:"city" validate:"required"`
	Address *string `json:"address" validate:"required"`
	Region  *string `json:"region" validate:"required"`
	Email   *string `json:"email" validate:"required"`
}

// Payment provides information about the payment
// @Description payment information
type Payment struct {
	ID           *string `json:"id" validate:"required"`
	Transaction  *string `json:"transaction" validate:"required"`
	RequestID    *string `json:"requestID"`
	Currency     *string `json:"currency" validate:"required"`
	Provider     *string `json:"provider" validate:"required"`
	Amount       *int    `json:"amount" validate:"required"`
	PaymentDT    *int    `json:"paymentDT" validate:"required"`
	Bank         *string `json:"bank" validate:"required"`
	DeliveryCost *int    `json:"deliveryCost" validate:"required"`
	GoodsTotal   *int    `json:"goodsTotal" validate:"required"`
	CustomFee    *int    `json:"customFee" validate:"required"`
}

// Item provides information about the product
// @Description Product information
type Item struct {
	ID          *string `json:"id" validate:"required"`
	OrderUID    *string `json:"orderUID" validate:"required"`
	ChrtID      *int    `json:"chrtID" validate:"required"`
	TrackNumber *string `json:"trackNumber" validate:"required"`
	Price       *int    `json:"price" validate:"required"`
	Rid         *string `json:"rid" validate:"required"`
	Name        *string `json:"name" validate:"required"`
	Sale        *int    `json:"sale" validate:"required"`
	Size        *string `json:"size" validate:"required"`
	TotalPrice  *int    `json:"totalPrice" validate:"required"`
	NmID        *int    `json:"nmID" validate:"required"`
	Brand       *string `json:"brand" validate:"required"`
	Status      *int    `json:"status" validate:"required"`
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
