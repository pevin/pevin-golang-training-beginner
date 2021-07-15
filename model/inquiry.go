package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Inquiry struct {
	Id            string           `json:"id"`
	TransactionId string           `json:"transaction_id" validate:"required"`
	PaymentCode   string           `json:"payment_code" validate:"required"`
	Amount        *decimal.Decimal `json:"amount" validate:"required"`
	CreatedAt     time.Time        `json:"-"`
	UpdatedAt     time.Time        `json:"-"`
}
