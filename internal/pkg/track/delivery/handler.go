package delivery

import (
	"encoding/json"
	"net/http"
	"no_homomorphism/pkg/logger"
	"strconv"

	"github.com/gorilla/mux"
	"no_homomorphism/internal/pkg/track"
)

type TrackHandler struct {
	TrackUC track.UseCase
	Log     *logger.MainLogger
}

func (h *TrackHandler) GetTrack(w http.ResponseWriter, r *http.Request) {
	varId, e := mux.Vars(r)["id"]
	if e == false {
		h.Log.HttpInfo("", "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ID, err := strconv.Atoi(varId)
	if err != nil {
		h.Log.HttpInfo("", "failed to parse id:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	track, err := h.TrackUC.GetTrackById(uint(ID))
	if err != nil {
		h.Log.HttpInfo("", "failed get track"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writer := json.NewEncoder(w)
	err = writer.Encode(track)
	if err != nil {
		h.Log.HttpInfo("", "can't write track info into json:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.Log.HttpInfo("", "OK", http.StatusOK)
}
