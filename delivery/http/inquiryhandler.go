package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pevin/pevin-golang-training-beginner/model"
	"github.com/pevin/pevin-golang-training-beginner/usecase"
	"gopkg.in/go-playground/validator.v9"
)

type InquiryHandler struct {
	Usecase usecase.IInquiryUseCase
}

func NewInquiryRouteHandler(usecase usecase.IInquiryUseCase) *InquiryHandler {
	return &InquiryHandler{
		Usecase: usecase,
	}
}

func (p *InquiryHandler) createInquiry(w http.ResponseWriter, r *http.Request) (err error) {
	inquiry, err := p.Usecase.InitFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validateError, err := p.validate(inquiry)
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

	err = p.Usecase.Create(r.Context(), &inquiry)
	if err != nil {
		if err.Error() == "Inquiry transaction already exist" {
			e := model.Error{Message: "Transaction ID already used."}
			resp, _ := json.Marshal(e)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(inquiry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

	return
}

func (p *InquiryHandler) validate(inquiry model.Inquiry) (valError model.Error, err error) {
	validate := validator.New()
	validateErrors := validate.Struct(inquiry)
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

func (p *InquiryHandler) InquiryRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		p.createInquiry(w, r)
		return
	default:
		NotFoundHandler(w, r)
	}
}
