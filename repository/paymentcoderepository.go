package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/pevin/pevin-golang-training-beginner/model"

	_ "github.com/lib/pq"
)

type IPaymentCodeRepository interface {
	Create(ctx context.Context, p *model.PaymentCode) (err error)
	Get(ctx context.Context, id string) (paymentCode model.PaymentCode, err error)
	GetIdsToBeExpired(ctx context.Context) (ids []string, err error)
	UpdateStatusById(ctx context.Context, id string, status string) (err error)
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

func (r PaymentCodeRepository) UpdateStatusById(ctx context.Context, id string, status string) (err error) {
	res, err := r.Db.ExecContext(
		context.Background(),
		"UPDATE payment_codes SET status=$2 WHERE id=$1",
		id, status,
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

// This method will fetch all payment code ids with ACTIVE status and due expired date
func (r PaymentCodeRepository) GetIdsToBeExpired(ctx context.Context) (ids []string, err error) {
	now := time.Now().UTC()
	rows, err := r.Db.QueryContext(context.Background(), "SELECT id FROM payment_codes where expiration_date <= $1 and status = $2", now, "ACTIVE")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}

	return
}
