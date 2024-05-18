package database

import (
	"context"
	"encoding/json"
	"os"
	"service/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// wrapper around pgxpool
type Database struct {
	pool Pool
}
type Pool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// On conflict of id's does nothing
func (db *Database) InsertOrder(order *models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, err = db.pool.Exec(context.Background(),
		`INSERT INTO orders (id, data, insert_time) 
        VALUES ($1, $2, $3)
        ON CONFLICT (id)
        DO NOTHING
        `,
		order.Id, data, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	return nil
}

// Returns models.Order on success, else error
func (db *Database) GetOrder(uid string) (*models.Order, error) {
	var data []byte
	err := db.pool.QueryRow(context.Background(),
		`SELECT data FROM orders WHERE id = $1`, uid).Scan(&data)
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
		`SELECT data FROM orders
        ORDER BY insert_time DESC 
        LIMIT $1`, maxAmount)
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
	url := os.Getenv("DB_URL")
	if len(url) == 0 {
		url = "localhost"
	}
	pool, err := pgxpool.New(context.Background(), "postgres://order_service:password@"+url+"/order_db")
	if err != nil {
		panic("Unable to connect to db")
	}
	return &Database{
		pool: pool,
	}
}
