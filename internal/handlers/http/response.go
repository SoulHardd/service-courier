package http

import (
	"avito/internal/handlers/http/httperror"
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if response != nil {
		json.NewEncoder(w).Encode(response)
	}
}

func WriteErrorResponse(w http.ResponseWriter, err interface{}) {
	if httpErr, ok := err.(httperror.HTTPError); ok {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}
	if domainErr, ok := err.(error); ok {
		httpErr := httperror.MapDomainError(domainErr)
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
}
