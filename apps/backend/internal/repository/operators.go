package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Kost0/L0/internal/models"
)

// SQLOrderRepository provides information about database
type SQLOrderRepository struct {
	DB *sql.DB
}

// NewOrderRepository create new SQLOrderRepository
// Accepts:
//   - db: database
//
// Returns:
//   - *SQLOrderRepository
func NewOrderRepository(db *sql.DB) *SQLOrderRepository {
	return &SQLOrderRepository{DB: db}
}

// InsertOrder insert data to database
// Accepts:
//   - data: all data about order
//
// Returns:
//   - error if something wrong
func (r *SQLOrderRepository) InsertOrder(data *models.CombinedData) error {
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
`

	queryOrder := `
INSERT INTO orders (
    order_uid,
    track_number,
    entry,
    delivery_id,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
`
	queryItem := `
INSERT INTO items (
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
`

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(queryDelivery,
		*delivery.ID,
		*delivery.Name,
		*delivery.Phone,
		*delivery.Zip,
		*delivery.City,
		*delivery.Address,
		*delivery.Region,
		*delivery.Email,
	)
	if err != nil {
		errRollBack := tx.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return err
	}

	_, err = tx.Exec(queryOrder,
		order.OrderUID,
		*order.TrackNumber,
		*order.Entry,
		*order.DeliveryID,
		*order.Locale,
		*order.InternalSignature,
		*order.CustomerID,
		*order.DeliveryService,
		*order.Shardkey,
		*order.SmID,
		*order.DateCreated,
		*order.OofShard,
	)
	if err != nil {
		errRollBack := tx.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return err
	}

	_, err = tx.Exec(queryPayment,
		*payment.Transaction,
		*payment.RequestID,
		*payment.Currency,
		*payment.Provider,
		*payment.Amount,
		*payment.PaymentDT,
		*payment.Bank,
		*payment.DeliveryCost,
		*payment.GoodsTotal,
		*payment.CustomFee,
	)
	if err != nil {
		errRollBack := tx.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return err
	}

	for _, i := range items {
		_, err = tx.Exec(queryItem,
			*i.ChrtID,
			*i.TrackNumber,
			*i.Price,
			*i.Rid,
			*i.Name,
			*i.Sale,
			*i.Size,
			*i.TotalPrice,
			*i.NmID,
			*i.Brand,
			*i.Status,
		)
		if err != nil {
			errRollBack := tx.Rollback()
			if errRollBack != nil {
				return errRollBack
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Println("Data inserted in db")
	return nil
}

// InsertWithRetry insert data to database using multiple attempts if necessary
// Accepts:
//   - ctx: context
//   - data: all data about order
//
// Returns:
//   - error if something wrong
func (r *SQLOrderRepository) InsertWithRetry(ctx context.Context, data *models.CombinedData) error {
	maxRetries := 5
	delay := time.Millisecond * 50
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := r.InsertOrder(data)
		if err == nil {
			return nil
		}

		if !(errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone)) {
			return err
		}

		delay *= 2

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("insert failed, retry after %d attempts", maxRetries)
}

// SelectOrder select data from database
// Accepts:
//   - orderID: identifier
//
// Returns:
//   - all data about order
//   - error if something wrong
func (r *SQLOrderRepository) SelectOrder(orderUID string) (data *models.CombinedData, err error) {
	order := models.Order{}
	delivery := models.Delivery{}
	payment := models.Payment{}
	items := []models.Item{}

	queryOrder := `SELECT * FROM orders WHERE order_uid = $1`
	row := r.DB.QueryRow(queryOrder, orderUID)
	err = row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.DeliveryID,
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
	row = r.DB.QueryRow(queryDelivery, &order.DeliveryID)
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

	queryPayment := `SELECT * FROM payment WHERE transaction = $1`
	row = r.DB.QueryRow(queryPayment, &order.OrderUID)
	err = row.Scan(
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

	queryItems := `SELECT * FROM items WHERE track_number = $1`
	rows, err := r.DB.Query(queryItems, &order.TrackNumber)
	if err != nil {
		return nil, err
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			err = errClose
		}
	}()

	for rows.Next() {
		item := models.Item{}
		err = rows.Scan(
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

	data = &models.CombinedData{
		Order:    order,
		Payment:  payment,
		Delivery: delivery,
		Items:    items,
	}
	return
}

// SelectWithRetry select data from database using multiple attempts if necessary
// Accepts:
//   - ctx: context
//   - orderID: identifier
//
// Returns:
//   - all data about order
//   - error if something wrong
func (r *SQLOrderRepository) SelectWithRetry(ctx context.Context, orderUID string) (*models.CombinedData, error) {
	maxRetries := 5
	delay := time.Millisecond * 50
	for attempt := 0; attempt < maxRetries; attempt++ {
		data, err := r.SelectOrder(orderUID)
		if err == nil {
			return data, nil
		}

		if !(errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone)) {
			return nil, err
		}

		delay *= 2

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("select failed, retry after %d attempts", maxRetries)
}
