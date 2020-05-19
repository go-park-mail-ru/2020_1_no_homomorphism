package delivery

import (
	"encoding/json"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/album"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/gorilla/mux"
	"net/http"
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

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		user = models.User{Id: ""}
	}

	albumData, err := h.AlbumUC.GetAlbumById(varId, user.Id)
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
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "album delivery", "GetUserAlbums", "failed to from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

func (h *AlbumHandler) RateAlbum(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "delivery", "RateAlbum", "failed to get from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	varId, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.AlbumUC.RateAlbum(varId, user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to rate albums"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
