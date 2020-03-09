package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"no_homomorphism/internal/pkg/track"
)

type TrackHandler struct {
	TrackUC track.UseCase
	Log     *logger.MainLogger
}

func (h *TrackHandler) GetTrack(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ID, err := strconv.Atoi(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to parse id:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	track, err := h.TrackUC.GetTrackById(uint(ID))
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed get track"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writer := json.NewEncoder(w)
	err = writer.Encode(track)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write track info into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
