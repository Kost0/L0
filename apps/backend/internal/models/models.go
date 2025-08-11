package models

import "time"

type Order struct {
	OrderUID          string    `json:"orderUID"`
	TrackNumber       string    `json:"trackNumber"`
	Entry             string    `json:"entry"`
	DeliveryID        string    `json:"deliveryID"`
	PaymentID         string    `json:"paymentID"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internalSignature"`
	CustomerID        string    `json:"customerID"`
	DeliveryService   string    `json:"deliveryService"`
	Shardkey          string    `json:"shardKey"`
	SmID              int       `json:"smID"`
	DateCreated       time.Time `json:"dateCreated"`
	OofShard          string    `json:"oofShard"`
}

type Delivery struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	ID           string `json:"id"`
	Transaction  string `json:"transaction"`
	RequestID    string `json:"requestID"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"paymentDT"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"deliveryCost"`
	GoodsTotal   int    `json:"goodsTotal"`
	CustomFee    int    `json:"customFee"`
}

type Item struct {
	ID          string `json:"id"`
	OrderUID    string `json:"orderUID"`
	ChrtID      int    `json:"chrtID"`
	TrackNumber string `json:"trackNumber"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"totalPrice"`
	NmID        int    `json:"nmID"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type CombinedData struct {
	Order    Order
	Payment  Payment
	Delivery Delivery
	Items    []Item
}
