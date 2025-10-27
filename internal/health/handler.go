package health

import (
	"net/http"

	"github.com/go-chi/render"
)

type CheckResponse struct {
	Status string `json:"status"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, &CheckResponse{Status: "UP"})
}
