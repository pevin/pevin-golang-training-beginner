package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/usecase"

	"gopkg.in/go-playground/validator.v9"

	_ "github.com/lib/pq"
)

type WebCommand struct{}

func (c *WebCommand) Execute() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/hello-world", helloWorldHandler)

	http.HandleFunc("/payment-codes", paymentCodeRouteHandler)
	http.HandleFunc("/payment-codes/", paymentCodeRouteHandler)

	http.HandleFunc("/", notFoundHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type PaymentCodeHandler struct {
	Usecase usecase.IPaymentCodeUseCase
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
		notFoundHandler(w, r)
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
func paymentCodeRouteHandler(w http.ResponseWriter, r *http.Request) {
	pcUsecase := initPaymentUsecase()
	pcHandler := PaymentCodeHandler{
		Usecase: pcUsecase,
	}

	switch r.Method {
	case "POST":
		pcHandler.createPaymentCode(w, r)
		return
	case "GET":
		pcHandler.getPaymentCodeHandler(w, r)
		return
	default:
		notFoundHandler(w, r)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	error := model.Error{Message: "Request not found!"}

	resp, err := json.Marshal(error)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(resp)
}
