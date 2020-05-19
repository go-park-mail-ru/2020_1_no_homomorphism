package delivery

import (
	"encoding/json"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/artist"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ArtistHandler struct {
	ArtistUC artist.UseCase
	TrackUC  track.UseCase
	Log      *logger.MainLogger
}

func (h *ArtistHandler) GetFullArtistInfo(w http.ResponseWriter, r *http.Request) {
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

	artistInfo, err := h.ArtistUC.GetArtistById(varId, user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "cant get artistInfo:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(&artistInfo)
	if err != nil {
		h.Log.LogWarning(r.Context(), "artist delivery", "GetFullArtistInfo", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *ArtistHandler) GetBoundedArtists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	start, okStart := vars["start"]
	end, okEnd := vars["end"]

	if !okStart || !okEnd {
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

	artists, err := h.ArtistUC.GetBoundedArtists(uStart, uEnd)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get artists"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(struct {
		Artists []models.Artist `json:"artists"`
	}{artists})

	if err != nil {
		h.Log.LogWarning(r.Context(), "artist delivery", "GetBoundedArtists", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *ArtistHandler) GetArtistStat(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no data in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	artistStat, err := h.ArtistUC.GetArtistStat(id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get artist's stat"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(artistStat)
	if err != nil {
		h.Log.LogWarning(r.Context(), "artist delivery", "GetArtistStat", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *ArtistHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "artist delivery", "Subscription", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no data in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.ArtistUC.Subscription(id, user.Id)
	if err != nil {
		h.Log.LogWarning(r.Context(), "artist delivery", "Subscription", "failed to subscribe on artist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *ArtistHandler) SubscriptionList(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		h.Log.LogWarning(r.Context(), "playlist delivery", "GetUserPlaylists", "failed to get from ctx")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subscriptions, err := h.ArtistUC.SubscriptionList(user.Id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get user's subscriptions"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(subscriptions)
	if err != nil {
		h.Log.LogWarning(r.Context(), "artist delivery", "SubscriptionList", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
