package usecase

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/repository"

	"github.com/google/uuid"
)

type IPaymentUseCase interface {
	InitFromRequest(r *http.Request) (inquiry model.Payment, err error)
	Create(ctx context.Context, p *model.Payment) (err error)
}
type PaymentUseCase struct {
	Repo repository.IPaymentRepository
}

func (u PaymentUseCase) InitFromRequest(r *http.Request) (payment model.Payment, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	json.Unmarshal([]byte(body), &payment)

	return
}

func (u PaymentUseCase) Create(ctx context.Context, payment *model.Payment) (err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	payment.Id = id.String()

	now := time.Now().UTC()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	err = u.Repo.Create(ctx, payment)

	if err != nil {
		return
	}

	return
}
