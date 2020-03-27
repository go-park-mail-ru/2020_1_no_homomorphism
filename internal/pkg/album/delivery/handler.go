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

func (h *AlbumHandler) GetAlbumTracks(w http.ResponseWriter, r *http.Request) {
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

	tracks, err := h.TrackUC.GetTracksByAlbumId(varId)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get album' tracks"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(struct {
		Album  models.Album   `json:"album"`
		Count  int            `json:"count"`
		Tracks []models.Track `json:"tracks"`
	}{albumData, len(tracks), tracks})

	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't write album' tracks into json:"+err.Error(), http.StatusBadRequest)
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
		h.Log.HttpInfo(r.Context(), "can't write album into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
