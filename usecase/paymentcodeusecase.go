package usecase

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pevin-golang-training-beginner/model"
	"pevin-golang-training-beginner/producer"
	"pevin-golang-training-beginner/repository"
	"time"

	"github.com/google/uuid"
)

type PaymentCodeUseCase struct {
	Repo     repository.IPaymentCodeRepository
	Producer producer.IPaymentCodeMessageProducer
}

func (u PaymentCodeUseCase) InitFromRequest(r *http.Request) (paymentCode model.PaymentCode, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	json.Unmarshal([]byte(body), &paymentCode)

	return
}

func (u PaymentCodeUseCase) Create(ctx context.Context, paymentCode *model.PaymentCode) (err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	paymentCode.Id = id.String()

	now := time.Now().UTC()
	paymentCode.CreatedAt = now
	paymentCode.UpdatedAt = now

	expDate := now.AddDate(50, 0, 0)
	paymentCode.ExpirationDate = expDate

	paymentCode.Status = model.PAYMENT_CODE_STATUS_ACTIVE

	err = u.Repo.Create(ctx, paymentCode)

	if err != nil {
		return
	}

	err = u.Producer.Produce(paymentCode)

	return
}

func (u PaymentCodeUseCase) Get(ctx context.Context, id string) (p model.PaymentCode, err error) {
	p, err = u.Repo.Get(ctx, id)
	if err != nil {
		return
	}

	err = u.Producer.Produce(&p)

	return
}
