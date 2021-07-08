package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/usecase"
	"gopkg.in/go-playground/validator.v9"
)

type PaymentCodeHandler struct {
	Usecase usecase.IPaymentCodeUseCase
}

func NewPaymentCodeRouteHandler(usecase usecase.IPaymentCodeUseCase) *PaymentCodeHandler {
	return &PaymentCodeHandler{
		Usecase: usecase,
	}

}

func (p *PaymentCodeHandler) createPaymentCode(w http.ResponseWriter, r *http.Request) (err error) {
	paymentCode, err := p.Usecase.InitFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validateError, err := p.validate(paymentCode)
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

	err = p.Usecase.Create(r.Context(), &paymentCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(paymentCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

	return
}

func (p *PaymentCodeHandler) getPaymentCodeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/payment-codes/")

	paymentCode, err := p.Usecase.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if paymentCode.Id == "" {
		NotFoundHandler(w, r)
		return
	}

	resp, _ := json.Marshal(paymentCode)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (p *PaymentCodeHandler) validate(paymentCode model.PaymentCode) (valError model.Error, err error) {
	validate := validator.New()
	validateErrors := validate.Struct(paymentCode)
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

// PAYMENT CODE HANDLERS
func (p *PaymentCodeHandler) PaymentCodeRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		p.createPaymentCode(w, r)
		return
	case "GET":
		p.getPaymentCodeHandler(w, r)
		return
	default:
		NotFoundHandler(w, r)
	}
}
