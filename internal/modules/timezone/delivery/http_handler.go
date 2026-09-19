package delivery

import (
	"encoding/json"
	"net/http"
	"strings"

	"viewgo/internal/modules/timezone/application"
)

type TimezoneHTTPHandler struct {
	useCase *application.TimezoneUseCase
}

func NewTimezoneHTTPHandler(useCase *application.TimezoneUseCase) *TimezoneHTTPHandler {
	return &TimezoneHTTPHandler{
		useCase: useCase,
	}
}

func (h *TimezoneHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/timezones", h.handleTimezones)
	mux.HandleFunc("/api/v1/timezones/", h.handleTimezoneByName)
}

func (h *TimezoneHTTPHandler) handleTimezones(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	ctx := r.Context()

	var timezones interface{}
	var err error

	if query != "" {
		timezones, err = h.useCase.Search(ctx, query)
	} else {
		timezones, err = h.useCase.ListAll(ctx)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   timezones,
	})
}

func (h *TimezoneHTTPHandler) handleTimezoneByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/api/v1/timezones/")
	if name == "" {
		h.handleTimezones(w, r)
		return
	}

	tz, err := h.useCase.GetByName(r.Context(), name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   tz,
	})
}
