package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_repository "github.com/pevin/pevin-golang-training-beginner/mock/repository"
	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/repository"
)

func TestPaymentUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	type fields struct {
		Repo repository.IPaymentRepository
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := PaymentUseCase{
				Repo: tt.fields.Repo,
			}
			if err := u.Create(tt.args.ctx, tt.args.payment); (err != nil) != tt.wantErr {
				t.Errorf("PaymentUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
