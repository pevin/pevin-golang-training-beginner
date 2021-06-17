package usecase

import (
	"context"
	"errors"
	mock_producer "pevin-golang-training-beginner/mock/producer"
	mock_repository "pevin-golang-training-beginner/mock/repository"
	"pevin-golang-training-beginner/model"
	"pevin-golang-training-beginner/producer"
	"pevin-golang-training-beginner/repository"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPaymentCodeUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	err := errors.New("Mock Error")

	type fields struct {
		Repo     repository.IPaymentCodeRepository
		Producer producer.IPaymentCodeMessageProducer
	}
	type args struct {
		ctx         context.Context
		paymentCode *model.PaymentCode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "get-success",
			fields: fields{
				Repo: func() repository.IPaymentCodeRepository {
					repo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
					repo.
						EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil)
					return repo
				}(),
				Producer: func() producer.IPaymentCodeMessageProducer {
					producer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)
					producer.
						EXPECT().
						Produce(gomock.Any()).
						Return(nil)
					return producer
				}(),
			},
			args: args{
				ctx: context.TODO(),
				paymentCode: &model.PaymentCode{
					PaymentCode: "test-payment-code",
					Name:        "test name",
					Status:      "test-status",
				},
			},
		},
		{
			name: "with-error-in-repo",
			fields: fields{
				Repo: func() repository.IPaymentCodeRepository {
					repo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
					repo.
						EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(err)
					return repo
				}(),
				Producer: func() producer.IPaymentCodeMessageProducer {
					producer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)
					return producer
				}(),
			},
			args: args{
				ctx: context.TODO(),
				paymentCode: &model.PaymentCode{
					PaymentCode: "test-payment-code",
					Name:        "test name",
					Status:      "test-status",
				},
			},
			wantErr: true,
		},
		{
			name: "with-error-in-producer",
			fields: fields{
				Repo: func() repository.IPaymentCodeRepository {
					repo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
					repo.
						EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil)
					return repo
				}(),
				Producer: func() producer.IPaymentCodeMessageProducer {
					producer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)
					producer.
						EXPECT().
						Produce(gomock.Any()).
						Return(err)
					return producer
				}(),
			},
			args: args{
				ctx: context.TODO(),
				paymentCode: &model.PaymentCode{
					PaymentCode: "test-payment-code",
					Name:        "test name",
					Status:      "test-status",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := PaymentCodeUseCase{
				Repo:     tt.fields.Repo,
				Producer: tt.fields.Producer,
			}
			if err := u.Create(tt.args.ctx, tt.args.paymentCode); (err != nil) != tt.wantErr {
				t.Errorf("PaymentCodeUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestPaymentCodeUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	paymentCode := model.PaymentCode{
		Id:          "test-id",
		PaymentCode: "test-payment-code",
		Name:        "test name",
		Status:      "test-status",
	}
	emptyPaymentCode := model.PaymentCode{}
	err := errors.New("Mock Error")

	type fields struct {
		Repo     repository.IPaymentCodeRepository
		Producer producer.IPaymentCodeMessageProducer
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantP   model.PaymentCode
		wantErr bool
	}{
		{
			name: "get-success",
			fields: fields{
				Repo: func() repository.IPaymentCodeRepository {
					repo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
					repo.
						EXPECT().
						Get(gomock.Any(), "test-id").
						Return(paymentCode, nil)
					return repo
				}(),
				Producer: func() producer.IPaymentCodeMessageProducer {
					producer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)
					producer.
						EXPECT().
						Produce(gomock.Any()).
						Return(nil)
					return producer
				}(),
			},
			args: args{
				ctx: context.TODO(),
				id:  "test-id",
			},
			wantP: paymentCode,
		},
		{
			name: "get-error-from-repo",
			fields: fields{
				Repo: func() repository.IPaymentCodeRepository {
					repo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
					repo.
						EXPECT().
						Get(gomock.Any(), "test-id").
						Return(emptyPaymentCode, err)
					return repo
				}(),
				Producer: func() producer.IPaymentCodeMessageProducer {
					producer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)
					return producer
				}(),
			},
			args: args{
				ctx: context.TODO(),
				id:  "test-id",
			},
			wantP:   emptyPaymentCode,
			wantErr: true,
		},
		{
			name: "get-error-from-producer",
			fields: fields{
				Repo: func() repository.IPaymentCodeRepository {
					repo := mock_repository.NewMockIPaymentCodeRepository(ctrl)
					repo.
						EXPECT().
						Get(gomock.Any(), "test-id").
						Return(paymentCode, nil)
					return repo
				}(),
				Producer: func() producer.IPaymentCodeMessageProducer {
					producer := mock_producer.NewMockIPaymentCodeMessageProducer(ctrl)
					producer.
						EXPECT().
						Produce(gomock.Any()).
						Return(err)
					return producer
				}(),
			},
			args: args{
				ctx: context.TODO(),
				id:  "test-id",
			},
			wantP:   paymentCode,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := PaymentCodeUseCase{
				Repo:     tt.fields.Repo,
				Producer: tt.fields.Producer,
			}
			gotP, err := u.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentCodeUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("PaymentCodeUseCase.Get() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}
