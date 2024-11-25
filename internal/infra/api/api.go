package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"app/internal/core/model"
)

type handlers interface {
	GetEmployeesByName(
		ctx context.Context,
		req model.GetEmployeesRequest,
	) ([]model.Employee, error)
}

func newHandler(h handlers) http.Handler {
	r := chi.NewRouter()

	r.Get("/employees/{name}", func(w http.ResponseWriter, r *http.Request) {
		req, err := parseGetEmployeesRequest(r)
		if err != nil {
			log.Printf("failed to parse the get employees request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		emps, err := h.GetEmployeesByName(r.Context(), req)
		if err != nil {
			log.Printf("failed to get employees by name: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(emps)
		if err != nil {
			log.Printf("failed to marshal the employees list: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write(body); err != nil {
			log.Printf("failed to write the employees list: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	return r
}

func parseGetEmployeesRequest(r *http.Request) (model.GetEmployeesRequest, error) {
	req := model.GetEmployeesRequest{
		Name: chi.URLParam(r, "name"),
	}
	countParam := r.URL.Query().Get("limit")
	lastIDParam := r.URL.Query().Get("last-id")

	var err error
	req.Limit, err = strconv.Atoi(countParam)
	if err != nil {
		return req, fmt.Errorf("failed to parse the limit parameter %s: %w", countParam, err)
	}
	req.LastID, err = strconv.Atoi(lastIDParam)
	if err != nil {
		return req, fmt.Errorf("failed to parse the last-id parameter %s: %w", lastIDParam, err)
	}

	return req, nil
}
