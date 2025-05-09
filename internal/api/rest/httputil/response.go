package httputil

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, payload interface{}) {
	RespondWithJSON(w, code, payload)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
