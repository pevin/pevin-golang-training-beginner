package model

import (
	"time"
)

type Payment struct {
	Id            string    `json:"id"`
	TransactionId string    `json:"transaction_id" validate:"required"`
	PaymentCode   string    `json:"payment_code" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Amount        string    `json:"amount" validate:"required"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}
