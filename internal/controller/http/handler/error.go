package handler

import (
	"encoding/json"
	"net/http"
)

func ErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"status": http.StatusText(code),
		"error":  error,
	})
}

func HandleError(w http.ResponseWriter, err error, statusCode int) {
	ErrorJSON(w, err.Error(), statusCode)
}

func DecodingError(w http.ResponseWriter) {
	ErrorJSON(w, "body decoding error", http.StatusBadRequest)
}

func QueryError(w http.ResponseWriter) {
	ErrorJSON(w, "query request not full", http.StatusBadRequest)
}
