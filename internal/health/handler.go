package health

import (
	"encoding/json"
	"net/http"
)

type CheckResponse struct {
	Status string `json:"status"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&CheckResponse{Status: "UP"})
}
