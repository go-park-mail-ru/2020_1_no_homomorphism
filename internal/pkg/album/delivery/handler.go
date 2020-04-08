package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"no_homomorphism/internal/pkg/album"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/pkg/logger"
)

type AlbumHandler struct {
	AlbumUC album.UseCase
	TrackUC track.UseCase
	Log     *logger.MainLogger
}

func (h *AlbumHandler) GetFullAlbum(w http.ResponseWriter, r *http.Request) {
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	albumData, err := h.AlbumUC.GetAlbumById(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get album data"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(albumData)

	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write album data into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *AlbumHandler) GetUserAlbums(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(models.User)

	albums, err := h.AlbumUC.GetUserAlbums(user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get user' albums"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(struct {
		Albums []models.Album `json:"albums"`
	}{albums})

	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write albums into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *AlbumHandler) GetBoundedAlbumsByArtistId(w http.ResponseWriter, r *http.Request) {
	artistId, okId := r.Context().Value("id").(string)
	start, okStart := r.Context().Value("start").(uint64)
	end, okEnd := r.Context().Value("end").(uint64)

	if !okId || !okStart || !okEnd {
		h.Log.LogWarning(r.Context(), "album delivery", "GetBoundedArtistTracks", "failed to get vars")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	albums, err := h.AlbumUC.GetBoundedAlbumsByArtistId(artistId, start, end)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get albums"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Id     string         `json:"id"`
		Albums []models.Album `json:"albums"`
	}{artistId, albums})

	if err != nil {
		h.Log.LogWarning(r.Context(), "album delivery", "GetBoundedAlbumsByArtistId", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *AlbumHandler) GetBoundedAlbumTracks(w http.ResponseWriter, r *http.Request) {
	id, okId := r.Context().Value("id").(string)
	start, okStart := r.Context().Value("start").(uint64)
	end, okEnd := r.Context().Value("end").(uint64)

	if !okId || !okStart || !okEnd {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetBoundedPlaylistTracks", "failed to get vars")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracks, err := h.TrackUC.GetBoundedTracksByAlbumId(id, start, end)
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
