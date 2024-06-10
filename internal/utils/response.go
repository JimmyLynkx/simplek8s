package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	response := map[string]interface{}{
		"msg":  "failure",
		"data": message,
	}
	json.NewEncoder(w).Encode(response)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	response := map[string]interface{}{
		"msg":  "success",
		"data": payload,
	}
	json.NewEncoder(w).Encode(response)
}
