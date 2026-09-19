package delivery

import (
	"encoding/json"
	"net/http"
	"strings"

	"viewgo/internal/modules/country/application"
)

type CountryHTTPHandler struct {
	useCase *application.CountryUseCase
}

func NewCountryHTTPHandler(useCase *application.CountryUseCase) *CountryHTTPHandler {
	return &CountryHTTPHandler{
		useCase: useCase,
	}
}

func (h *CountryHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/countries", h.handleCountries)
	mux.HandleFunc("/api/v1/countries/", h.handleCountryByCode)
}

func (h *CountryHTTPHandler) handleCountries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	ctx := r.Context()

	var countries interface{}
	var err error

	if query != "" {
		countries, err = h.useCase.Search(ctx, query)
	} else {
		countries, err = h.useCase.ListAll(ctx)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   countries,
	})
}

func (h *CountryHTTPHandler) handleCountryByCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/api/v1/countries/")
	if code == "" {
		h.handleCountries(w, r)
		return
	}

	country, err := h.useCase.GetByCode(r.Context(), code)
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
		"data":   country,
	})
}
