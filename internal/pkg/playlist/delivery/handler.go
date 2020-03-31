package delivery

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/playlist"
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
)

type PlaylistHandler struct {
	PlaylistUC playlist.UseCase
	TrackUC    track.UseCase
	Log        *logger.MainLogger
}

func (h *PlaylistHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(models.User)

	playlists, err := h.PlaylistUC.GetUserPlaylists(user.Id)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to get playlists"+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(struct {
		Playlists []models.Playlist `json:"playlists"`
	}{playlists})

	if err != nil {
		h.sendBadRequest(w, r.Context(), "can't write playlists info into json:"+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) GetFullPlaylistById(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(models.User)
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no id in mux vars")
		return
	}
	ok, err := h.PlaylistUC.CheckAccessToPlaylist(user.Id, varId)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to check access: "+err.Error())
		return
	}
	if !ok {
		h.Log.HttpInfo(r.Context(), "current user doesnt have access to the playlist", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	playlistData, err := h.PlaylistUC.GetPlaylistById(varId)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to get playlistData: "+err.Error())
		return
	}

	tracks, err := h.TrackUC.GetTracksByPlaylistId(varId)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to get tracks: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(struct {
		Playlist models.Playlist `json:"playlist"`
		Count    int             `json:"count"`
		Tracks   []models.Track  `json:"tracks"`
	}{playlistData, len(tracks), tracks})

	if err != nil {
		h.sendBadRequest(w, r.Context(), "can't write tracks info into json: "+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) sendBadRequest(w http.ResponseWriter, ctx context.Context, msg string) {
	h.Log.HttpInfo(ctx, msg, http.StatusBadRequest)
	w.WriteHeader(http.StatusBadRequest)
}
