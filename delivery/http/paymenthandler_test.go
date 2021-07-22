package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	mock_http "github.com/pevin/pevin-golang-training-beginner/mock/net/http"
	mock_usecase "github.com/pevin/pevin-golang-training-beginner/mock/usecase"
	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/usecase"
)

func TestWebPaymentHandler_createPayment(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockErr := errors.New("Mock Error")

	type fields struct {
		InquiryUsecase usecase.IInquiryUseCase
		PaymentUsecase usecase.IPaymentUseCase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	inq := model.Inquiry{
		TransactionId: "test-transaction-id",
		PaymentCode:   "test-payment-code",
		Amount:        "100",
		Id:            "test-id",
	}
	payment := model.Payment{
		TransactionId: "test-transaction-id",
		PaymentCode:   "test-payment-code",
		Name:          "test-name",
		Amount:        "100",
	}

	j, err := json.Marshal(inq)

	req, err := http.NewRequest("POST", "/inquiries", strings.NewReader(string(j)))
	if err != nil {
		t.Fatal(err)
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
				InquiryUsecase: func() usecase.IInquiryUseCase {
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
					uc.
						EXPECT().
						GetByTransactionId(gomock.Any(), inq.TransactionId).
						Return(inq, nil)
					return uc
				}(),
				PaymentUsecase: func() usecase.IPaymentUseCase {
					uc := mock_usecase.NewMockIPaymentUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(payment, nil)
					uc.
						EXPECT().
						Create(gomock.Any(), &payment).
						Return(nil)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().WriteHeader(http.StatusCreated)

					resp, _ := json.Marshal(payment)
					rw.EXPECT().Write(resp).Return(0, nil)

					return rw
				}(),
				r: req,
			},
		},
		{
			name: "get-bad-request-for-invalid-input",
			fields: fields{
				InquiryUsecase: func() usecase.IInquiryUseCase {
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
					return uc
				}(),
				PaymentUsecase: func() usecase.IPaymentUseCase {
					invalidPc := payment
					invalidPc.TransactionId = ""
					invalidPc.PaymentCode = ""
					uc := mock_usecase.NewMockIPaymentUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(invalidPc, nil)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().WriteHeader(http.StatusBadRequest)

					rw.EXPECT().Write(gomock.Any()).Return(0, nil)

					return rw
				}(),
				r: req,
			},
		},
		{
			name: "get-bad-request-for-trx-not-exist",
			fields: fields{
				InquiryUsecase: func() usecase.IInquiryUseCase {
					emptyInq := model.Inquiry{}
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
					uc.
						EXPECT().
						GetByTransactionId(gomock.Any(), payment.TransactionId).
						Return(emptyInq, nil)
					return uc
				}(),
				PaymentUsecase: func() usecase.IPaymentUseCase {
					uc := mock_usecase.NewMockIPaymentUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(payment, nil)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().WriteHeader(http.StatusBadRequest)

					rw.EXPECT().Write(gomock.Any()).Return(0, nil)

					return rw
				}(),
				r: req,
			},
		},
		{
			name: "get-internal-error-from-usecase",
			fields: fields{
				InquiryUsecase: func() usecase.IInquiryUseCase {
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
					uc.
						EXPECT().
						GetByTransactionId(gomock.Any(), inq.TransactionId).
						Return(inq, nil)
					return uc
				}(),
				PaymentUsecase: func() usecase.IPaymentUseCase {
					invalidP := payment
					invalidP.TransactionId = ""
					invalidP.PaymentCode = ""
					uc := mock_usecase.NewMockIPaymentUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(payment, nil)
					uc.
						EXPECT().
						Create(gomock.Any(), &payment).
						Return(mockErr)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().WriteHeader(http.StatusInternalServerError)

					rw.EXPECT().Write(gomock.Any()).Return(0, nil)

					return rw
				}(),
				r: req,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaymentHandler{
				InquiryUsecase: tt.fields.InquiryUsecase,
				PaymentUsecase: tt.fields.PaymentUsecase,
			}
			if err := p.createPayment(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("PaymentHandler.createPayment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
