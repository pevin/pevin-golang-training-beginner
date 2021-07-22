package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/usecase"
	"gopkg.in/go-playground/validator.v9"
)

type PaymentHandler struct {
	InquiryUsecase usecase.IInquiryUseCase
	PaymentUsecase usecase.IPaymentUseCase
}

func NewPaymentRouteHandler(
	inquiryUsecase usecase.IInquiryUseCase,
	paymentUsecase usecase.IPaymentUseCase,
) *PaymentHandler {
	return &PaymentHandler{
		InquiryUsecase: inquiryUsecase,
		PaymentUsecase: paymentUsecase,
	}
}

func (p *PaymentHandler) createPayment(w http.ResponseWriter, r *http.Request) (err error) {
	payment, err := p.PaymentUsecase.InitFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validateError, err := p.validate(payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if validateError.Message != "" {
		resp, _ := json.Marshal(validateError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	inquiry, err := p.InquiryUsecase.GetByTransactionId(r.Context(), payment.TransactionId)
	if len(inquiry.Id) == 0 {
		e := model.Error{Message: "Transaction ID not found."}
		resp, _ := json.Marshal(e)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	err = p.PaymentUsecase.Create(r.Context(), &payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(payment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

	return
}

func (p *PaymentHandler) validate(payment model.Payment) (valError model.Error, err error) {
	validate := validator.New()
	validateErrors := validate.Struct(payment)
	if validateErrors != nil {
		for _, validateErrors := range validateErrors.(validator.ValidationErrors) {
			switch validateErrors.Tag() {
			case "required":
				valError = model.Error{Message: fmt.Sprintf("field '%s' is required", validateErrors.Field())}
				return
			}
		}
	}
	return
}

func (p *PaymentHandler) PaymentRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		p.createPayment(w, r)
		return
	default:
		NotFoundHandler(w, r)
	}
}
