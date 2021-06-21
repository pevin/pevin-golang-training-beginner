package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pevin-golang-training-beginner/model"
	"pevin-golang-training-beginner/usecase"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

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

// PAYMENT CODE HANDLERS
func paymentCodeRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createPaymentCodeHandler(w, r)
		return
	case "GET":
		getPaymentCodeHandler(w, r)
		return
	default:
		notFoundHandler(w, r)
	}
}

func createPaymentCodeHandler(w http.ResponseWriter, r *http.Request) {
	usecase := usecase.PaymentCodeUseCase{}
	paymentCode, err := usecase.InitFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validateError, err := validate(paymentCode)
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

	err = usecase.Create(r.Context(), &paymentCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(paymentCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func getPaymentCodeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/payment-codes/")

	usecase := usecase.PaymentCodeUseCase{}
	paymentCode, err := usecase.Get(r.Context(), id)
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

func validate(paymentCode model.PaymentCode) (valError model.Error, err error) {
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

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/hello-world", helloWorldHandler)

	http.HandleFunc("/payment-codes", paymentCodeRouteHandler)
	http.HandleFunc("/payment-codes/", paymentCodeRouteHandler)

	http.HandleFunc("/", notFoundHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
