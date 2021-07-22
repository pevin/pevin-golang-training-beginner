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

func TestWebInquiryHandler_createInquiry(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockErr := errors.New("Mock Error")
	trxExistErr := errors.New("Inquiry transaction already exist")

	type fields struct {
		Usecase usecase.IInquiryUseCase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	inq := model.Inquiry{
		TransactionId: "test-transaction-id",
		PaymentCode:   "test-payment-code",
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
				Usecase: func() usecase.IInquiryUseCase {
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(inq, nil)
					uc.
						EXPECT().
						Create(gomock.Any(), &inq).
						Return(nil)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().WriteHeader(http.StatusCreated)

					resp, _ := json.Marshal(inq)
					rw.EXPECT().Write(resp).Return(0, nil)

					return rw
				}(),
				r: req,
			},
		},
		{
			name: "get-bad-request-for-invalid-input",
			fields: fields{
				Usecase: func() usecase.IInquiryUseCase {
					invalidPc := inq
					invalidPc.TransactionId = ""
					invalidPc.PaymentCode = ""
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
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
			name: "get-bad-request-for-trx-already-exist",
			fields: fields{
				Usecase: func() usecase.IInquiryUseCase {
					invalidPc := inq
					invalidPc.TransactionId = ""
					invalidPc.PaymentCode = ""
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(inq, nil)
					uc.
						EXPECT().
						Create(gomock.Any(), &inq).
						Return(trxExistErr)
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
			wantErr: true,
		},
		{
			name: "get-internal-error-from-usecase",
			fields: fields{
				Usecase: func() usecase.IInquiryUseCase {
					invalidPc := inq
					invalidPc.TransactionId = ""
					invalidPc.PaymentCode = ""
					uc := mock_usecase.NewMockIInquiryUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(inq, nil)
					uc.
						EXPECT().
						Create(gomock.Any(), &inq).
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
			p := &InquiryHandler{
				Usecase: tt.fields.Usecase,
			}
			if err := p.createInquiry(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("InquiryHandler.createInquiry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
