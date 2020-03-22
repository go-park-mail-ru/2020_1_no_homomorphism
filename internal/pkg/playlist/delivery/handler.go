package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/playlist"
	"no_homomorphism/pkg/logger"
)

type PlaylistHandler struct {
	PlaylistUC playlist.UseCase
	Log        *logger.MainLogger
}

func (h *PlaylistHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(*models.User)

	playlists, err := h.PlaylistUC.GetUserPlaylists(user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get playlists"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(playlists)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write playlists info into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) GetPlaylistTracks(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playlists, err := h.PlaylistUC.GetPlaylistWithTracks(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get playlists: "+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(playlists)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write playlists info into json: "+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
