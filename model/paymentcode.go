package model

import (
	"time"
)

const (
	PAYMENT_CODE_STATUS_ACTIVE   = "ACTIVE"
	PAYMENT_CODE_STATUS_INACTIVE = "INACTIVE"
	PAYMENT_CODE_STATUS_EXPIRED  = "EXPIRED"
)

type PaymentCode struct {
	Id             string    `json:"id"`
	PaymentCode    string    `json:"payment_code" validate:"required"`
	Name           string    `json:"name" validate:"required"`
	Status         string    `json:"status"`
	ExpirationDate time.Time `json:"expiration_date"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
