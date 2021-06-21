package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"pevin-golang-training-beginner/model"

	_ "github.com/lib/pq"
)

type PaymentCodeRepository struct {
	db *sql.DB
}

func (r PaymentCodeRepository) Create(ctx context.Context, p *model.PaymentCode) (err error) {
	res, err := r.getDB().ExecContext(
		context.Background(),
		"INSERT INTO payment_codes (id, payment_code, name, status, expiration_date, created_at, updated_at) VALUES($1 ,$2 ,$3, $4, $5, $6, $7)",
		p.Id, p.PaymentCode, p.Name, p.Status, p.ExpirationDate, p.CreatedAt, p.UpdatedAt,
	)

	if err != nil {
		log.Fatal(err)
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
	rows, err := r.getDB().QueryContext(context.Background(), "SELECT id, payment_code, name, status FROM payment_codes where id = $1 limit 1", id)
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

func (r PaymentCodeRepository) getDB() *sql.DB {
	if r.db == nil {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")

		pgDsn := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName)

		db, err := sql.Open("postgres", pgDsn)
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}
		r.db = db
	}

	return r.db
}
