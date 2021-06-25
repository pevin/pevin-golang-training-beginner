package usecase

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/producer"
	"github.com/pevin/pevin-golang-training-beginner/repository"

	"github.com/google/uuid"
)

type IPaymentCodeUseCase interface {
	InitFromRequest(r *http.Request) (paymentCode model.PaymentCode, err error)
	Create(ctx context.Context, p *model.PaymentCode) (err error)
	Get(ctx context.Context, id string) (paymentCode model.PaymentCode, err error)
	ExpireWithPassDueExpiryDate(ctx context.Context) (err error)
}
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

func (u PaymentCodeUseCase) ExpireWithPassDueExpiryDate(ctx context.Context) (err error) {
	r, err := u.Repo.GetIdsToBeExpired(ctx)
	if err != nil {
		return
	}
	for _, id := range r {
		err = u.Repo.UpdateStatusById(ctx, id, "EXPIRED")
		if err != nil {
			return
		}
	}

	return
}
