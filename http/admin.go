package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/thoas/stats"
)

type adminHandler struct {
	Stats *stats.Stats
}

func newAdminHandler(s *stats.Stats) http.Handler {
	return handlers.MethodHandler{
		"GET": &adminHandler{
			Stats: s,
		},
	}
}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := h.Stats.Data()

	b, _ := json.Marshal(stats)

	w.Write(b)
}
