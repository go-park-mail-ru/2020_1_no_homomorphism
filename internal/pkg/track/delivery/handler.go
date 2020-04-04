package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	track "no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
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
	trackData, err := h.TrackUC.GetTrackById(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed get trackData"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writer := json.NewEncoder(w)
	err = writer.Encode(trackData)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write trackData info into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

