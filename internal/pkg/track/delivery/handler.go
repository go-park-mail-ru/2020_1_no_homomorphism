package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"no_homomorphism/internal/pkg/models"
	track "no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
	"strconv"
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

func (h *TrackHandler) GetAlbumTracks(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlists, err := h.TrackUC.GetTracksByAlbumId(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get album' track"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(playlists)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write album' track into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *TrackHandler) GetBoundedArtistTracks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistId, okId := vars["id"]
	start, okStart := vars["start"]
	end, okEnd := vars["end"]

	if !okId || !okStart || !okEnd {
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

	tracks, err := h.TrackUC.GetBoundedTracksByArtistId(artistId, uStart, uEnd)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get tracks"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Id     string         `json:"id"`
		Tracks []models.Track `json:"tracks"`
	}{artistId, tracks})

	if err != nil {
		h.Log.LogWarning(r.Context(), "tracks delivery", "GetBoundedArtistTracks", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
