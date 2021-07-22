package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/pevin/pevin-golang-training-beginner/model"

	_ "github.com/lib/pq"
)

type IInquiryRepository interface {
	Create(ctx context.Context, p *model.Inquiry) (err error)
	GetByTransactionId(ctx context.Context, trxId string) (inquiry model.Inquiry, err error)
}

type InquiryRepository struct {
	Db *sql.DB
}

func (r InquiryRepository) Create(ctx context.Context, p *model.Inquiry) (err error) {
	res, err := r.Db.ExecContext(
		context.Background(),
		"INSERT INTO inquiries (id, transaction_id, payment_code, amount, created_at, updated_at) VALUES($1 ,$2 ,$3, $4, $5, $6)",
		p.Id, p.TransactionId, p.PaymentCode, p.Amount, p.CreatedAt, p.UpdatedAt,
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

func (r InquiryRepository) GetByTransactionId(ctx context.Context, trxId string) (inquiry model.Inquiry, err error) {
	rows, err := r.Db.QueryContext(context.Background(), "SELECT id, transaction_id, payment_code, amount FROM inquiries where transaction_id = $1 limit 1", trxId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&inquiry.Id,
			&inquiry.TransactionId,
			&inquiry.PaymentCode,
			&inquiry.Amount,
		); err != nil {
			log.Fatal(err)
		}
		return
	}

	return
}
