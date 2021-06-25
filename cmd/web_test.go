package cmd

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

func TestWebPaymentCodeHandler_createPaymentCode(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockErr := errors.New("Mock Error")

	type fields struct {
		Usecase usecase.IPaymentCodeUseCase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	pc := model.PaymentCode{
		Name:        "test-name",
		PaymentCode: "test-payment-code",
	}

	j, err := json.Marshal(pc)

	req, err := http.NewRequest("POST", "/payment-codes", strings.NewReader(string(j)))
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
				Usecase: func() usecase.IPaymentCodeUseCase {
					uc := mock_usecase.NewMockIPaymentCodeUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(pc, nil)
					uc.
						EXPECT().
						Create(gomock.Any(), &pc).
						Return(nil)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().WriteHeader(http.StatusCreated)

					resp, _ := json.Marshal(pc)
					rw.EXPECT().Write(resp).Return(0, nil)

					return rw
				}(),
				r: req,
			},
		},
		{
			name: "get-bad-request-for-invalid-input",
			fields: fields{
				Usecase: func() usecase.IPaymentCodeUseCase {
					invalidPc := pc
					invalidPc.Name = ""
					invalidPc.PaymentCode = ""
					uc := mock_usecase.NewMockIPaymentCodeUseCase(ctrl)
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
			name: "get-internal-error-from-usecase",
			fields: fields{
				Usecase: func() usecase.IPaymentCodeUseCase {
					invalidPc := pc
					invalidPc.Name = ""
					invalidPc.PaymentCode = ""
					uc := mock_usecase.NewMockIPaymentCodeUseCase(ctrl)
					uc.
						EXPECT().
						InitFromRequest(gomock.Any()).
						Return(pc, nil)
					uc.
						EXPECT().
						Create(gomock.Any(), &pc).
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
			p := &PaymentCodeHandler{
				Usecase: tt.fields.Usecase,
			}
			if err := p.createPaymentCode(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("PaymentCodeHandler.createPaymentCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebPaymentCodeHandler_getPaymentCodeHandler(t *testing.T) {
	type fields struct {
		Usecase usecase.IPaymentCodeUseCase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockErr := errors.New("Mock Error")

	pc := model.PaymentCode{
		Id:          "test-id",
		Name:        "test-name",
		PaymentCode: "test-payment-code",
	}

	req, err := http.NewRequest("GET", "/payment-codes/test-id", nil)

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "get-success",
			fields: fields{
				Usecase: func() usecase.IPaymentCodeUseCase {
					uc := mock_usecase.NewMockIPaymentCodeUseCase(ctrl)
					uc.
						EXPECT().
						Get(gomock.Any(), pc.Id).
						Return(pc, nil)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)

					resp, _ := json.Marshal(pc)
					rw.EXPECT().Write(resp).Return(0, nil)

					return rw
				}(),
				r: req,
			},
		},
		{
			name: "get-not-found",
			fields: fields{
				Usecase: func() usecase.IPaymentCodeUseCase {
					uc := mock_usecase.NewMockIPaymentCodeUseCase(ctrl)
					uc.
						EXPECT().
						Get(gomock.Any(), pc.Id).
						Return(model.PaymentCode{}, nil)
					return uc
				}(),
			},
			args: args{
				w: func() http.ResponseWriter {
					rw := mock_http.NewMockResponseWriter(ctrl)

					rw.EXPECT().Header().Return(req.Header)
					rw.EXPECT().WriteHeader(http.StatusNotFound)
					error := model.Error{Message: "Request not found!"}
					resp, _ := json.Marshal(error)
					rw.EXPECT().Write(resp).Return(0, nil)

					return rw
				}(),
				r: req,
			},
		},
		{
			name: "get-internal-server-error-from-usecase",
			fields: fields{
				Usecase: func() usecase.IPaymentCodeUseCase {
					uc := mock_usecase.NewMockIPaymentCodeUseCase(ctrl)
					uc.
						EXPECT().
						Get(gomock.Any(), pc.Id).
						Return(pc, mockErr)
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaymentCodeHandler{
				Usecase: tt.fields.Usecase,
			}
			p.getPaymentCodeHandler(tt.args.w, tt.args.r)
		})
	}

}
