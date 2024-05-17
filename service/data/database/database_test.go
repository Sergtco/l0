package database

import (
	"encoding/json"
	"service/models"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

type AnyTime struct{}

func (a AnyTime) Match(v interface{}) bool {
	_, ok := v.(int64)
	return ok
}

func TestInsert(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	order := &models.Order{Id: "Hello"}
	data, _ := json.Marshal(order)
	mock.ExpectExec("INSERT INTO orders").
		WithArgs(order.Id, data, AnyTime{}).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	db := Database{pool: mock}
	if err := db.InsertOrder(order); err != nil {
		t.Errorf("Insertion failed: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations weren't met %s", err)
		return
	}
}

func TestSelecting(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	order := &models.Order{Id: "Test"}
	data, _ := json.Marshal(order)
	row := pgxmock.NewRows([]string{"data"}).AddRow(data)

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(order.Id, data, AnyTime{}).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectQuery("SELECT data").
		WithArgs(order.Id).
		WillReturnRows(row)

	db := Database{pool: mock}
	if err := db.InsertOrder(order); err != nil {
		t.Errorf("Insertion failed: %s", err)
		return
	}

	if _, err := db.GetOrder(order.Id); err != nil {
		t.Errorf("Unable to get %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations weren't met %s", err)
	}
}
