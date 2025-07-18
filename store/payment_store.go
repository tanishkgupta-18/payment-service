package store

import (
	"gofr.dev/pkg/gofr"
)

type Payment struct {
	ID     int
	Amount float64
	Status string
}

type PaymentStore struct{}

func NewPaymentStore() *PaymentStore {
	return &PaymentStore{}
}

func (s *PaymentStore) CreatePayment(ctx *gofr.Context, p *Payment) (int, error) {
	query := `INSERT INTO payments (amount, status) VALUES (?, ?)`
	res, err := ctx.SQL.ExecContext(ctx, query, p.Amount, p.Status)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (s *PaymentStore) UpdatePaymentStatus(ctx *gofr.Context, id int, status string) error {
	query := `UPDATE payments SET status = ? WHERE id = ?`
	_, err := ctx.SQL.ExecContext(ctx, query, status, id)
	return err
}
