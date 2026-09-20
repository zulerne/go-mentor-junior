package common

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondJSON(log *slog.Logger, w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("failed to encode response", "error", err)
	}
}
