package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/pevin/pevin-golang-training-beginner/model"

	_ "github.com/lib/pq"
)

type IPaymentRepository interface {
	Create(ctx context.Context, p *model.Payment) (err error)
}

type PaymentRepository struct {
	Db *sql.DB
}

func (r PaymentRepository) Create(ctx context.Context, p *model.Payment) (err error) {
	res, err := r.Db.ExecContext(
		context.Background(),
		"INSERT INTO payments (id, transaction_id, payment_code, name, amount, created_at, updated_at) VALUES($1 ,$2 ,$3, $4, $5, $6, $7)",
		p.Id, p.TransactionId, p.PaymentCode, p.Name, p.Amount, p.CreatedAt, p.UpdatedAt,
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	rowAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatal(err)
		return
	}

	if rowAffected != 1 {
		log.Fatalf("Expected row affected equal to 1 but got %d", rowAffected)
		return
	}

	return
}
