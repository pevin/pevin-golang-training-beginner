package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pevin/pevin-golang-training-beginner/producer"
	"github.com/pevin/pevin-golang-training-beginner/repository"
	"github.com/pevin/pevin-golang-training-beginner/usecase"
)

func initPaymentUsecase() usecase.IPaymentCodeUseCase {
	pcRepo := repository.PaymentCodeRepository{Db: initDb()}
	pcProducer := producer.PaymentCodeMessageProducer{}
	pcUsecase := usecase.PaymentCodeUseCase{Repo: pcRepo, Producer: pcProducer}

	return pcUsecase
}

func initDb() *sql.DB {
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "postgres")
	dbName := getEnv("DB_NAME", "traingolang")

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
	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
