package repository

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/Kost0/L0/internal/models"
	"github.com/segmentio/kafka-go"
)

func InsertOrder(db *sql.DB, msg *kafka.Message) error {
	var data models.CombinedData
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		return err
	}

	delivery := data.Delivery
	payment := data.Payment
	order := data.Order
	items := data.Items

	queryDelivery := `
INSERT INTO delivery (
    id,
    name,
    phone,
    zip,
    city,
    address,
    region,
    email
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
`

	queryPayment := `
INSERT INTO payment (
    id,
    transaction,
    request_id,
    currency,
    provider,
    amount,
    payment_dt,
    bank,
    delivery_cost,
    goods_total,
    custom_fee
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

	queryOrder := `
INSERT INTO orders (
    order_uid,
    track_number,
    entry,
    delivery_id,
    payment_id,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);
`

	queryItem := `
INSERT INTO items (
    id,
    order_id,
    chrt_id,
    track_number,
    price,
    rid,
    name,
    sale,
    size,
    total_price,
    nm_id,
    brand,
    status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`

	_, err = db.Exec(queryDelivery,
		delivery.ID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)
	if err != nil {
		return err
	}
	log.Println("Delivery inserted")

	_, err = db.Exec(queryPayment,
		payment.ID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)
	if err != nil {
		return err
	}
	log.Println("payment inserted")

	_, err = db.Exec(queryOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.DeliveryID,
		order.PaymentID,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return err
	}
	log.Println("Order inserted")

	for _, i := range items {
		_, err = db.Exec(queryItem,
			i.ID,
			i.OrderUID,
			i.ChrtID,
			i.TrackNumber,
			i.Price,
			i.Rid,
			i.Name,
			i.Sale,
			i.Size,
			i.TotalPrice,
			i.NmID,
			i.Brand,
			i.Status,
		)
		if err != nil {
			return err
		}
	}
	log.Println("Items inserted")
	return nil
}

func SelectOrder(db *sql.DB, orderUID string) (*models.CombinedData, error) {
	order := models.Order{}
	delivery := models.Delivery{}
	payment := models.Payment{}
	items := []models.Item{}

	queryOrder := `SELECT * FROM orders WHERE order_uid = $1`
	row := db.QueryRow(queryOrder, orderUID)
	err := row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.DeliveryID,
		&order.PaymentID,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, err
	}

	queryDelivery := `SELECT * FROM delivery WHERE id = $1`
	row = db.QueryRow(queryDelivery, order.DeliveryID)
	err = row.Scan(
		&delivery.ID,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		return nil, err
	}

	queryPayment := `SELECT * FROM payment WHERE id = $1`
	row = db.QueryRow(queryPayment, order.PaymentID)
	err = row.Scan(
		&payment.ID,
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}

	queryItems := `SELECT * FROM items WHERE order_id = $1`
	rows, err := db.Query(queryItems, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Item{}
		err = rows.Scan(
			&item.ID,
			&item.OrderUID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	data := models.CombinedData{
		Order:    order,
		Payment:  payment,
		Delivery: delivery,
		Items:    items,
	}
	return &data, nil
}
