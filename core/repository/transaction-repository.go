package repository

import (
	"database/sql"
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
	"strings"
)

type TransactionRepositoryInterface interface {
	GetDb() *sql.DB

	// order
	InsertOrder(order *entity.Order) error
	UpdateOrder(tx *sql.Tx, req *entity.UpdateOrderRequest) error
	GetOrderById(id string, isActive *bool) (*entity.Order, error)

	// order_item
	InsertOrderItems(items []*entity.OrderItem) error
	GetOrderItemsByOrderId(orderId string) ([]*entity.OrderItem, error)

	// payment
	InsertPayment(tx *sql.Tx, req *entity.Payment) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepositoryInterface {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) GetDb() *sql.DB {
	return r.db
}

func (r *transactionRepository) InsertOrder(order *entity.Order) error {
	query := `INSERT INTO orders (id, user_id, shop_id, amount, status, created_at, expired_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query, order.Id, order.UserId, order.ShopId, order.Amount, order.Status, order.CreatedAt, order.ExpiredAt)
	if err != nil {
		return fmt.Errorf("error repo insert order: %v", err.Error())
	}

	return nil
}

func (r *transactionRepository) UpdateOrder(tx *sql.Tx, req *entity.UpdateOrderRequest) error {
	query := `UPDATE orders 
				SET status = $1 
				WHERE id = $2`

	_, err := tx.Exec(query, req.Status, req.OrderId)
	if err != nil {
		return fmt.Errorf("error repo update order: %v", err.Error())
	}

	return nil
}

func (r *transactionRepository) GetOrderById(id string, isActive *bool) (*entity.Order, error) {
	query := `SELECT id, user_id, shop_id, status, amount, created_at, expired_at 
				FROM orders 
				WHERE id = $1`

	if isActive != nil {
		if *isActive {
			query += " AND expired_at >= NOW()"
		} else {
			query += " AND expired_at < NOW()"
		}
	}

	order := &entity.Order{}

	err := r.db.QueryRow(query, id).Scan(&order.Id, &order.UserId, &order.ShopId, &order.Status, &order.Amount, &order.CreatedAt, &order.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorutil.NewErrorCode(errorutil.ErrNotFound, fmt.Errorf("error repo get order: order id '%v' is not found", id))
		} else {
			return nil, fmt.Errorf("error repo get order: %v", err.Error())
		}
	}

	return order, nil
}

func (r *transactionRepository) InsertOrderItems(items []*entity.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, shop_id, warehouse_id, quantity, unit_price) VALUES %s`

	values := []interface{}{}
	placeholders := []string{}

	for i, item := range items {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*3+1, i*3+2, i*3+3, i*3+4, i*3+5, i*3+6))
		values = append(values, item.OrderId, item.ProductId, item.ShopId, item.WarehouseId, item.Quantity, item.UnitPrice)
	}

	var queryValues string
	queryValues += strings.Join(placeholders, ", ")

	query = fmt.Sprintf(query, queryValues)

	_, err := r.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("error repo insert order items: %v", err.Error())
	}

	return nil
}

func (r *transactionRepository) GetOrderItemsByOrderId(orderId string) ([]*entity.OrderItem, error) {
	query := `SELECT order_id, product_id, shop_id, warehouse_id, quantity, unit_price 
				FROM order_items 
				WHERE order_id = $1`

	rows, err := r.db.Query(query, orderId)
	if err != nil {
		return nil, fmt.Errorf("error repo get order items: %v", err.Error())
	}

	var items []*entity.OrderItem
	for rows.Next() {
		item := &entity.OrderItem{}
		err := rows.Scan(&item.OrderId, &item.ProductId, &item.ShopId, &item.WarehouseId, &item.Quantity, &item.UnitPrice)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *transactionRepository) InsertPayment(tx *sql.Tx, req *entity.Payment) error {
	query := `INSERT INTO payments (id, order_id, amount, status) 
				VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(query, req.Id, req.OrderId, req.Amount, req.Status)
	if err != nil {
		return fmt.Errorf("error repo insert payment: %v", err.Error())
	}

	return nil
}
