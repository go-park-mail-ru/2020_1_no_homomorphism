package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	tracks "no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
)

type TrackHandler struct {
	TrackUC tracks.UseCase
	Log     *logger.MainLogger
}

func (h *TrackHandler) GetTrack(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	track, err := h.TrackUC.GetTrackById(varId)
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

func (h *TrackHandler) GetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlists, err := h.TrackUC.GetTracksByAlbumId(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get album' tracks"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(playlists)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write album' tracks into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

