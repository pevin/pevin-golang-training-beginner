package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/repository"

	"github.com/google/uuid"
)

type IInquiryUseCase interface {
	InitFromRequest(r *http.Request) (inquiry model.Inquiry, err error)
	Create(ctx context.Context, p *model.Inquiry) (err error)
	GetByTransactionId(ctx context.Context, id string) (inquiry model.Inquiry, err error)
}
type InquiryUseCase struct {
	Repo repository.IInquiryRepository
}

func (u InquiryUseCase) InitFromRequest(r *http.Request) (inquiry model.Inquiry, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	json.Unmarshal([]byte(body), &inquiry)

	return
}

func (u InquiryUseCase) Create(ctx context.Context, inquiry *model.Inquiry) (err error) {
	existingTrx, err := u.GetByTransactionId(ctx, inquiry.TransactionId)
	if err != nil {
		return
	}
	if len(existingTrx.Id) > 0 {
		err = errors.New("Inquiry transaction already exist")
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	inquiry.Id = id.String()

	now := time.Now().UTC()
	inquiry.CreatedAt = now
	inquiry.UpdatedAt = now

	err = u.Repo.Create(ctx, inquiry)

	if err != nil {
		return
	}

	return
}

func (u InquiryUseCase) GetByTransactionId(ctx context.Context, id string) (p model.Inquiry, err error) {
	p, err = u.Repo.GetByTransactionId(ctx, id)
	if err != nil {
		return
	}

	return
}
