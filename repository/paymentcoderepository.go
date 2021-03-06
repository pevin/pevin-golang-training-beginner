package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/pevin/pevin-golang-training-beginner/model"

	_ "github.com/lib/pq"
)

type IPaymentCodeRepository interface {
	Create(ctx context.Context, p *model.PaymentCode) (err error)
	Get(ctx context.Context, id string) (paymentCode model.PaymentCode, err error)
}

type PaymentCodeRepository struct {
	Db *sql.DB
}

func (r PaymentCodeRepository) Create(ctx context.Context, p *model.PaymentCode) (err error) {
	res, err := r.Db.ExecContext(
		context.Background(),
		"INSERT INTO payment_codes (id, payment_code, name, status, expiration_date, created_at, updated_at) VALUES($1 ,$2 ,$3, $4, $5, $6, $7)",
		p.Id, p.PaymentCode, p.Name, p.Status, p.ExpirationDate, p.CreatedAt, p.UpdatedAt,
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

func (r PaymentCodeRepository) Get(ctx context.Context, id string) (paymentCode model.PaymentCode, err error) {
	rows, err := r.Db.QueryContext(context.Background(), "SELECT id, payment_code, name, status FROM payment_codes where id = $1 limit 1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&paymentCode.Id,
			&paymentCode.PaymentCode,
			&paymentCode.Name,
			&paymentCode.Status,
		); err != nil {
			log.Fatal(err)
		}
		return
	}

	return
}
