package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/pkg/logger"
)

type PlaylistHandler struct {
	PlaylistUC playlist.UseCase
	TrackUC    track.UseCase
	Log        *logger.MainLogger
}

func (h *PlaylistHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetUserPlaylists", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no id in mux vars")
		return
	}

	if err := h.checkUserAccess(w, r, varId); err != nil {
		return
	}

	playlistData, err := h.PlaylistUC.GetPlaylistById(varId)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to get playlistData: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(playlistData)

	if err != nil {
		h.sendBadRequest(w, r.Context(), "can't write tracks info into json: "+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) GetBoundedPlaylistTracks(w http.ResponseWriter, r *http.Request) {
	id, okId := r.Context().Value("id").(string)
	start, okStart := r.Context().Value("start").(uint64)
	end, okEnd := r.Context().Value("end").(uint64)

	if !okId || !okStart || !okEnd {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetBoundedPlaylistTracks", "failed to get vars")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.checkUserAccess(w, r, id); err != nil {
		return
	}
	tracks, err := h.TrackUC.GetBoundedTracksByPlaylistId(id, start, end)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to get tracks"+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Id     string         `json:"id"`
		Tracks []models.Track `json:"tracks"`
	}{id, tracks})

	if err != nil {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetBoundedPlaylistTracks", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) sendBadRequest(w http.ResponseWriter, ctx context.Context, msg string) {
	h.Log.HttpInfo(ctx, msg, http.StatusBadRequest)
	w.WriteHeader(http.StatusBadRequest)
}

func (h *PlaylistHandler) checkUserAccess(w http.ResponseWriter, r *http.Request, varId string) error {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "checkUserAccess", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("failed to get from ctx")
	}

	ok, err := h.PlaylistUC.CheckAccessToPlaylist(user.Id, varId)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to check access: "+err.Error())
		return errors.New("failed to check access")
	}
	if !ok {
		h.Log.HttpInfo(r.Context(), "current user doesnt have access to the playlist", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return errors.New("no access")
	}
	return nil
}
