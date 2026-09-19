package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/currency/application"
)

type CurrencyHTTPHandler struct {
	useCase *application.CurrencyUseCase
}

func NewCurrencyHTTPHandler(useCase *application.CurrencyUseCase) *CurrencyHTTPHandler {
	return &CurrencyHTTPHandler{
		useCase: useCase,
	}
}

func (h *CurrencyHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/currencies", h.handleCurrencies)
	mux.HandleFunc("/api/v1/currencies/format", h.handleFormatCurrency)
	mux.HandleFunc("/api/v1/currencies/", h.handleCurrencyByCode)
}

func (h *CurrencyHTTPHandler) handleCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	ctx := r.Context()

	var currencies interface{}
	var err error

	if query != "" {
		currencies, err = h.useCase.Search(ctx, query)
	} else {
		currencies, err = h.useCase.ListAll(ctx)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   currencies,
	})
}

func (h *CurrencyHTTPHandler) handleCurrencyByCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/api/v1/currencies/")
	if code == "" || code == "format" {
		h.handleCurrencies(w, r)
		return
	}

	c, err := h.useCase.GetByCode(r.Context(), code)
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
		"data":   c,
	})
}

func (h *CurrencyHTTPHandler) handleFormatCurrency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")
	amountStr := r.URL.Query().Get("amount")

	if code == "" || amountStr == "" {
		http.Error(w, "Missing code or amount parameter", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "Invalid amount parameter", http.StatusBadRequest)
		return
	}

	formatted, err := h.useCase.FormatAmount(r.Context(), code, amount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"currency":  code,
		"amount":    amount,
		"formatted": formatted,
	})
}
