package usecase

import (
	"context"
	"errors"
	"testing"

	mock_producer "pevin-golang-training-beginner/mock/producer"
	mock_repository "pevin-golang-training-beginner/mock/repository"
	"pevin-golang-training-beginner/model"

	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mRepo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
	mProducer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)

	model := model.PaymentCode{
		PaymentCode: "test-payment-code",
		Name:        "test name",
		Status:      "test-status",
	}
	mRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)
	mProducer.
		EXPECT().
		Produce(gomock.Any()).
		Return(nil)

	usecase := PaymentCodeUseCase{Repo: mRepo, Producer: mProducer}

	got := usecase.Create(context.TODO(), &model)

	if got != nil {
		t.Errorf("got %q, wanted nil", got)
	}
}

func TestCreateWithErrorInRepo(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mRepo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
	mProducer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)

	model := model.PaymentCode{
		PaymentCode: "test-payment-code",
		Name:        "test name",
		Status:      "test-status",
	}

	err := errors.New("Mock Error")
	mRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).Return(err)

	usecase := PaymentCodeUseCase{Repo: mRepo, Producer: mProducer}

	got := usecase.Create(context.TODO(), &model)

	if got == nil {
		t.Errorf("got %q, wanted error", got)
	}
}
func TestCreateWithErrorInProducer(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mRepo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
	mProducer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)

	model := model.PaymentCode{
		PaymentCode: "test-payment-code",
		Name:        "test name",
		Status:      "test-status",
	}

	err := errors.New("Mock Error")
	mRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)
	mProducer.
		EXPECT().
		Produce(gomock.Any()).
		Return(err)

	usecase := PaymentCodeUseCase{Repo: mRepo, Producer: mProducer}

	got := usecase.Create(context.TODO(), &model)

	if got == nil {
		t.Errorf("got %q, wanted error", got)
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mRepo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
	mProducer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)

	model := model.PaymentCode{
		Id:          "test-id",
		PaymentCode: "test-payment-code",
		Name:        "test name",
		Status:      "test-status",
	}
	mRepo.
		EXPECT().
		Get(gomock.Any(), "test-id").
		Return(model, nil)
	mProducer.
		EXPECT().
		Produce(gomock.Any()).
		Return(nil)

	usecase := PaymentCodeUseCase{Repo: mRepo, Producer: mProducer}

	got, err := usecase.Get(context.TODO(), "test-id")

	if err != nil {
		t.Errorf("got %q, wanted nil", got)
	}

	if got != model {
		t.Errorf("got %q, wanted %q", got, model)
	}
}
func TestGetWithErrorInRepo(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mRepo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
	mProducer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)

	model := model.PaymentCode{}
	mockErr := errors.New("Mock Error")
	mRepo.
		EXPECT().
		Get(gomock.Any(), "test-id").
		Return(model, mockErr)

	usecase := PaymentCodeUseCase{Repo: mRepo, Producer: mProducer}

	got, err := usecase.Get(context.TODO(), "test-id")

	if err == nil {
		t.Errorf("got %q, wanted nil", got)
	}
}

func TestGetWithErrorInProducer(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mRepo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
	mProducer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)

	model := model.PaymentCode{}
	mockErr := errors.New("Mock Error")
	mRepo.
		EXPECT().
		Get(gomock.Any(), "test-id").
		Return(model, nil)
	mProducer.
		EXPECT().
		Produce(gomock.Any()).
		Return(mockErr)

	usecase := PaymentCodeUseCase{Repo: mRepo, Producer: mProducer}

	got, err := usecase.Get(context.TODO(), "test-id")

	if err == nil {
		t.Errorf("got %q, wanted nil", got)
	}
}
