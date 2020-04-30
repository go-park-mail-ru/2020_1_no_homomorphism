package delivery

import (
	"encoding/json"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/search"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type SearchHandler struct {
	SearchUC search.UseCase
	Log      *logger.MainLogger
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	varText, ok := mux.Vars(r)["text"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no text in vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	varCount, ok := mux.Vars(r)["count"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no count in vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	count, err := strconv.ParseUint(varCount, 10, 32)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to parse count", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchResult, err := h.SearchUC.Search(varText, uint(count))
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to search", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(searchResult)

	if err != nil {
		h.Log.LogWarning(r.Context(), "search delivery", "Search", "failed to encode: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
