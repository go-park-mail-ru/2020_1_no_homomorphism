package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"no_homomorphism/internal/pkg/track"
)

type TrackHandler struct {
	TrackUC track.UseCase
}
func (h *TrackHandler) GetTrack(w http.ResponseWriter, r *http.Request){
	varId, e := mux.Vars(r)["id"]
	if e == false {
		log.Println("no song id in mux vars")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ID, err := strconv.Atoi(varId)
	if err != nil {
		log.Println("song id is not u integer")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	track, err := h.TrackUC.GetTrackById(uint(ID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writer := json.NewEncoder(w)
	err = writer.Encode(track)
	if err != nil {
		log.Println("can't write track info into json :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
