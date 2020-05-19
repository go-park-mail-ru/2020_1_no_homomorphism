package delivery

import (
	"encoding/json"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	track "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/gorilla/mux"
	"net/http"
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

func (h *TrackHandler) GetBoundedArtistTracks(w http.ResponseWriter, r *http.Request) {
	id, okId := r.Context().Value("id").(string)
	start, okStart := r.Context().Value("start").(uint64)
	end, okEnd := r.Context().Value("end").(uint64)

	if !okId || !okStart || !okEnd {
		h.Log.LogWarning(r.Context(), "track delivery", "GetBoundedArtistTracks", "failed to get vars")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		user = models.User{Id: ""}
	}

	tracks, err := h.TrackUC.GetBoundedTracksByArtistId(id, start, end, user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get tracks"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Id     string         `json:"id"`
		Tracks []models.Track `json:"tracks"`
	}{id, tracks})

	if err != nil {
		h.Log.LogWarning(r.Context(), "track delivery", "GetBoundedArtistTracks", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *TrackHandler) GetBoundedAlbumTracks(w http.ResponseWriter, r *http.Request) {
	id, okId := r.Context().Value("id").(string)
	start, okStart := r.Context().Value("start").(uint64)
	end, okEnd := r.Context().Value("end").(uint64)

	if !okId || !okStart || !okEnd {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetBoundedPlaylistTracks", "failed to get vars")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		user = models.User{Id: ""}
	}

	tracks, err := h.TrackUC.GetBoundedTracksByAlbumId(id, start, end, user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get tracks"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Id     string         `json:"id"`
		Tracks []models.Track `json:"tracks"`
	}{id, tracks})

	if err != nil {
		h.Log.LogWarning(r.Context(), "tracks delivery", "GetBoundedAlbumTracks", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *TrackHandler) GetUserTracks(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "delivery", "GetUserTracks", "failed to get from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracks, err := h.TrackUC.GetUserTracks(user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get tracks"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Tracks []models.Track `json:"tracks"`
	}{tracks})

	if err != nil {
		h.Log.LogWarning(r.Context(), "tracks delivery", "GetBoundedAlbumTracks", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *TrackHandler) RateTrack(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "delivery", "GetUserTracks", "failed to get from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.TrackUC.RateTrack(user.Id, varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get tracks"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
