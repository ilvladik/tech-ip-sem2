package httpx

import (
	"encoding/json"
	"net/http"
)

func RespondJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondJson(w, status, map[string]string{"error": message})
}
