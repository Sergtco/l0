package database

import (
	"context"
	"encoding/json"
	"service/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func (db *Database) InsertOrder(order *models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	db.pool.Exec(context.Background(),
		`insert into orders (id, data) 
        values ($1, $2)
        on confilct (id)
        do nothing
        `,
		order.Id, data)
	return nil
}

func (db *Database) GetOrder(uid string) (*models.Order, error) {
	var data []byte
	err := db.pool.QueryRow(context.Background(),
		`select data from orders where id = $1`, uid).Scan(&data)
	if err != nil {
		return nil, err
	}
	out := &models.Order{}
	err = json.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (db *Database) GetTopOrders(maxAmount int) ([]*models.Order, error) {
	rows, err := db.pool.Query(context.Background(),
		`select data from orders limit $1`, maxAmount)
	if err != nil {
		return nil, err
	}
	rowsData, err := pgx.CollectRows(rows, pgx.RowTo[[]byte])
	if err != nil {
		return nil, err
	}
	orders := []*models.Order{}
	for _, row := range rowsData {
		order := &models.Order{}
		err := json.Unmarshal(row, order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func New() *Database {
	pool, err := pgxpool.New(context.Background(), "postgres://order_service:password@localhost/order_db")
	if err != nil {
		panic("Unable to connect to db")
	}
	return &Database{
		pool: pool,
	}
}
