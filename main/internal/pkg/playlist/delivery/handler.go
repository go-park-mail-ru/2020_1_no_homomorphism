package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/gorilla/mux"
	"net/http"
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

func (h *PlaylistHandler) checkUserAccess(w http.ResponseWriter, r *http.Request, playlistID string) error {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "checkUserAccess", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return errors.New("failed to get from ctx")
	}

	ok, err := h.PlaylistUC.CheckAccessToPlaylist(user.Id, playlistID)
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

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	name, ok := mux.Vars(r)["name"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no name in mux vars")
		return
	}

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "CreatePlaylist", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	plID, err := h.PlaylistUC.CreatePlaylist(name, user.Id)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	writer := json.NewEncoder(w)
	err = writer.Encode(struct {
		ID string `json:"playlist_id"`
	}{plID})

	if err != nil {
		h.sendBadRequest(w, r.Context(), "can't write into json: "+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusCreated)
}

func (h *PlaylistHandler) AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	plTracks := models.PlaylistTracks{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&plTracks)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "error while unmarshalling JSON:"+err.Error())
		return
	}

	if err = h.checkUserAccess(w, r, plTracks.PlaylistID); err != nil {
		return
	}

	if err = h.PlaylistUC.AddTrackToPlaylist(plTracks); err != nil {
		h.sendBadRequest(w, r.Context(), "cant add track to playlist:"+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) GetPlaylistsIDByTrack(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no id in mux vars")
		return
	}

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetPlaylistsIDByTrack", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	playlistIDs, err := h.PlaylistUC.GetPlaylistsIdByTrack(user.Id, id)

	if err != nil {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetPlaylistsIDByTrack", "failed to get playlists:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(struct {
		IDs []string `json:"playlists"`
	}{playlistIDs})

	if err != nil {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetPlaylistsIDByTrack", "failed to encode:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) DeleteTrackFromPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, ok := mux.Vars(r)["playlist"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no playlistID in mux vars")
		return
	}

	trackID, ok := mux.Vars(r)["track"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no trackID in mux vars")
		return
	}

	if err := h.checkUserAccess(w, r, playlistID); err != nil {
		return
	}

	if err := h.PlaylistUC.DeleteTrackFromPlaylist(playlistID, trackID); err != nil {
		h.sendBadRequest(w, r.Context(), "cant delete track from playlist:"+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, ok := mux.Vars(r)["id"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no id in mux vars")
		return
	}

	if err := h.checkUserAccess(w, r, playlistID); err != nil {
		return
	}

	if err := h.PlaylistUC.DeletePlaylist(playlistID); err != nil {
		h.sendBadRequest(w, r.Context(), "cant delete track from playlist:"+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
