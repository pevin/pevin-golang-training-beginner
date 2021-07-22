package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pevin/pevin-golang-training-beginner/delivery/sqs"
	mock_sqs "github.com/pevin/pevin-golang-training-beginner/mock/delivery/sqs"
	mock_repository "github.com/pevin/pevin-golang-training-beginner/mock/repository"
	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/repository"
)

func TestPaymentUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	err := errors.New("Mock Error")

	type fields struct {
		Repo      repository.IPaymentRepository
		Publisher sqs.IPaymentPublisher
	}
	type args struct {
		ctx     context.Context
		payment *model.Payment
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
				Repo: func() repository.IPaymentRepository {
					repo := mock_repository.NewMockIPaymentRepository(ctrl)
					repo.
						EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil)
					return repo
				}(),
				Publisher: func() sqs.IPaymentPublisher {
					publisher := mock_sqs.NewMockIPaymentPublisher(ctrl)
					publisher.
						EXPECT().
						Publish(gomock.Any()).
						Return(nil)
					return publisher
				}(),
			},
			args: args{
				ctx: context.TODO(),
				payment: &model.Payment{
					PaymentCode:   "test-payment-code",
					Name:          "test-name",
					TransactionId: "test-transaction-id",
					Amount:        "100000",
				},
			},
		},
		{
			name: "get-err-from-repo",
			fields: fields{
				Repo: func() repository.IPaymentRepository {
					repo := mock_repository.NewMockIPaymentRepository(ctrl)
					repo.
						EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(err)
					return repo
				}(),
				Publisher: func() sqs.IPaymentPublisher {
					publisher := mock_sqs.NewMockIPaymentPublisher(ctrl)
					return publisher
				}(),
			},
			args: args{
				ctx: context.TODO(),
				payment: &model.Payment{
					PaymentCode:   "test-payment-code",
					Name:          "test-name",
					TransactionId: "test-transaction-id",
					Amount:        "100000",
				},
			},
			wantErr: true,
		},
		{
			name: "get-err-from-publisher",
			fields: fields{
				Repo: func() repository.IPaymentRepository {
					repo := mock_repository.NewMockIPaymentRepository(ctrl)
					repo.
						EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil)
					return repo
				}(),
				Publisher: func() sqs.IPaymentPublisher {
					publisher := mock_sqs.NewMockIPaymentPublisher(ctrl)
					publisher.
						EXPECT().
						Publish(gomock.Any()).
						Return(err)
					return publisher
				}(),
			},
			args: args{
				ctx: context.TODO(),
				payment: &model.Payment{
					PaymentCode:   "test-payment-code",
					Name:          "test-name",
					TransactionId: "test-transaction-id",
					Amount:        "100000",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := PaymentUseCase{
				Repo:      tt.fields.Repo,
				Publisher: tt.fields.Publisher,
			}
			if err := u.Create(tt.args.ctx, tt.args.payment); (err != nil) != tt.wantErr {
				t.Errorf("PaymentUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
