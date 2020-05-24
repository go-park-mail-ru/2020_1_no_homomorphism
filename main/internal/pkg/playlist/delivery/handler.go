package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/middleware"
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
	ImgTypes   map[string]string
}

func (h *PlaylistHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(models.User)
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

	if err := h.checkUserAccess(w, r, varId, false); err != nil {
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
	id, okId := r.Context().Value(middleware.Id).(string)
	start, okStart := r.Context().Value(middleware.Start).(uint64)
	end, okEnd := r.Context().Value(middleware.End).(uint64)

	if !okId || !okStart || !okEnd {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetBoundedPlaylistTracks", "failed to get vars")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.checkUserAccess(w, r, id, false); err != nil {
		return
	}

	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		user = models.User{Id: ""}
	}

	tracks, err := h.TrackUC.GetBoundedTracksByPlaylistId(id, start, end, user.Id)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to get tracks"+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	output := models.PlaylistTracksArray{Id: id, Tracks: tracks}

	err = json.NewEncoder(w).Encode(output)

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

func (h *PlaylistHandler) checkUserAccess(w http.ResponseWriter, r *http.Request, playlistID string, isStrict bool) error {
	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		user = models.User{Id: "0"}
	}

	ok, err := h.PlaylistUC.CheckAccessToPlaylist(user.Id, playlistID, isStrict)
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

	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "CreatePlaylist", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	plID, err := h.PlaylistUC.CreatePlaylist(name, user.Id)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "cant create playlist:"+err.Error())
		return
	}

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

	if err = h.checkUserAccess(w, r, plTracks.PlaylistID, true); err != nil {
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

	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetPlaylistsIDByTrack", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	playlistIDs, err := h.PlaylistUC.GetUserPlaylistsIdByTrack(user.Id, id)

	if err != nil {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetPlaylistsIDByTrack", "failed to get playlists:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(models.PlaylistsID{IDs: playlistIDs})

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

	if err := h.checkUserAccess(w, r, playlistID, true); err != nil {
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

	if err := h.checkUserAccess(w, r, playlistID, true); err != nil {
		return
	}

	if err := h.PlaylistUC.DeletePlaylist(playlistID); err != nil {
		h.sendBadRequest(w, r.Context(), "cant delete track from playlist:"+err.Error())
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) ChangePrivacy(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no id in mux vars")
		return
	}

	err := h.PlaylistUC.ChangePrivacy(varId)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to change playlist privacy: "+err.Error())
		return
	}
}

func (h *PlaylistHandler) AddSharedPlaylist(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no id in mux vars")
		return
	}

	user, ok := r.Context().Value(middleware.UserKey).(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "AddSharedPlaylist", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.checkUserAccess(w, r, varId, false); err != nil {
		return
	}

	id, err := h.PlaylistUC.AddSharedPlaylist(varId, user.Id)
	if err != nil {
		h.sendBadRequest(w, r.Context(), "failed to copy playlist"+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	//output := models.Playlist{Id: id}

	err = json.NewEncoder(w).Encode(struct {
		Id string `json:"id"`
	}{id})

	if err != nil {
		h.Log.LogWarning(r.Context(), "playlist delivery", "AddSharedPlaylist", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *PlaylistHandler) UpdatePlaylistAvatar(w http.ResponseWriter, r *http.Request) {
	playlistID, ok := mux.Vars(r)["id"]
	if !ok {
		h.sendBadRequest(w, r.Context(), "no id in mux vars")
		return
	}
	//token, ok := r.Context().Value(middleware.CSRFTokenCorrect).(bool)
	//if !token || !ok {
	//	h.Log.HttpInfo(r.Context(), "permission denied: user has wrong csrf token", http.StatusUnauthorized)
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	if err := h.checkUserAccess(w, r, playlistID, true); err != nil {
		return
	}
	file, handler, err := r.FormFile("playlist_image")
	if err != nil || handler.Size == 0 {
		h.Log.HttpInfo(r.Context(), "can't read playlist_image", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := handler.Header.Get("Content-Type")
	elem, ok := h.ImgTypes[mimeType]
	if !ok {
		h.Log.HttpInfo(r.Context(), "wrong file content-type", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path, err := h.PlaylistUC.UpdateAvatar(playlistID, file, elem)
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "UpdatePlaylistAvatar", "failed to update avatar:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.Info("new file created:", path)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

func (h *PlaylistHandler) Update(w http.ResponseWriter, r *http.Request) {
	plID, okID := mux.Vars(r)["id"]
	plName, okName := mux.Vars(r)["name"]
	if !okID || !okName {
		h.sendBadRequest(w, r.Context(), "no name or id in mux vars")
		return
	}
	//token, ok := r.Context().Value(middleware.CSRFTokenCorrect).(bool)
	//if !token || !ok {
	//	h.Log.HttpInfo(r.Context(), "permission denied: user has wrong csrf token", http.StatusUnauthorized)
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	if err := h.checkUserAccess(w, r, plID, true); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.PlaylistUC.Update(plID, plName)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't update playlist:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
	w.WriteHeader(http.StatusOK)
}
