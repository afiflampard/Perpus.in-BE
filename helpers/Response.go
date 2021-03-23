package helpers

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Kode    uint
	Message string
	status  bool
}

func MessageResponses(status bool, kode uint, message string) *MessageResponse {
	return &MessageResponse{
		Kode:    kode,
		Message: message,
		status:  status,
	}
}
func ResponseWithError(w http.ResponseWriter, code int, message interface{}) {
	ResponseWithJson(w, code, map[string]interface{}{"error": message})
}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
