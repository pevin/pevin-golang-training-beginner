package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pevin/pevin-golang-training-beginner/model"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
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
