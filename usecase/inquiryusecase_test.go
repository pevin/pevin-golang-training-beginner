package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_repository "github.com/pevin/pevin-golang-training-beginner/mock/repository"
	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/repository"
)

func TestInquiryUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	type fields struct {
		Repo repository.IInquiryRepository
	}
	type args struct {
		ctx     context.Context
		inquiry *model.Inquiry
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
				Repo: func() repository.IInquiryRepository {
					repo := mock_repository.NewMockIInquiryRepository(ctrl)
					emptyInquiry := model.Inquiry{}
					repo.
						EXPECT().
						GetByTransactionId(gomock.Any(), gomock.Any()).
						Return(emptyInquiry, nil)
					repo.
						EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil)
					return repo
				}(),
			},
			args: args{
				ctx: context.TODO(),
				inquiry: &model.Inquiry{
					PaymentCode:   "test-payment-code",
					TransactionId: "test-transaction-id",
				},
			},
		},
		{
			name: "with-error-on-existing-transaction-id",
			fields: fields{
				Repo: func() repository.IInquiryRepository {
					repo := mock_repository.NewMockIInquiryRepository(ctrl)
					emptyInquiry := model.Inquiry{Id: "test-id"}
					repo.
						EXPECT().
						GetByTransactionId(gomock.Any(), gomock.Any()).
						Return(emptyInquiry, nil)
					return repo
				}(),
			},
			args: args{
				ctx: context.TODO(),
				inquiry: &model.Inquiry{
					PaymentCode:   "test-payment-code",
					TransactionId: "test-transaction-id",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := InquiryUseCase{
				Repo: tt.fields.Repo,
			}
			if err := u.Create(tt.args.ctx, tt.args.inquiry); (err != nil) != tt.wantErr {
				t.Errorf("InquiryUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
