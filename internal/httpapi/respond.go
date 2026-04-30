package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("err encoding to json")
	}
}

func writeError(w http.ResponseWriter, status int, code, message, field string) {
	writeJSON(w, status, errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
			Field:   field,
		},
	})
}

func badRequest(w http.ResponseWriter, code, msg, field string) {
	writeError(w, http.StatusBadRequest, code, msg, field)
}

func validationError(w http.ResponseWriter, code, msg, field string) {
	writeError(w, http.StatusBadRequest, code, msg, field)
}

func notFound(w http.ResponseWriter, code, msg, field string) {
	writeError(w, http.StatusNotFound, code, msg, field)
}

func conflict(w http.ResponseWriter, code, msg, field string) {
	writeError(w, http.StatusConflict, code, msg, field)
}

func internalError(w http.ResponseWriter, code, msg, field string) {
	writeError(w, http.StatusInternalServerError, code, msg, field)
}
