package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"no_homomorphism/internal/pkg/artist"
	"no_homomorphism/internal/pkg/csrf"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/pkg/logger"
	"strconv"
)

type ArtistHandler struct {
	ArtistUC        artist.UseCase
	SessionDelivery session.Delivery
	CSRF            csrf.CryptToken
	Log             *logger.MainLogger
}

func (h *ArtistHandler) GetFullArtistInfo(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	artistInfo, err := h.ArtistUC.GetArtistById(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "cant get artistInfo:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(&artistInfo)
	if err != nil {
		h.Log.LogWarning(r.Context(), "artistInfo delivery", "GetFullArtistInfo", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *ArtistHandler) GetBoundedArtists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	start, okStart := vars["start"]
	end, okEnd := vars["end"]

	if !okStart || !okEnd {
		h.Log.HttpInfo(r.Context(), "no data in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uStart, err1 := strconv.ParseUint(start, 10, 32)
	uEnd, err2 := strconv.ParseUint(end, 10, 32)
	if err1 != nil || err2 != nil {
		h.Log.HttpInfo(r.Context(), "failed to parse start or end parameters", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	artists, err := h.ArtistUC.GetBoundedArtists(uStart, uEnd)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get artists"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Artists []models.Artist `json:"artists"`
	}{artists})

	if err != nil {
		h.Log.LogWarning(r.Context(), "artist delivery", "GetBoundedArtists", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
