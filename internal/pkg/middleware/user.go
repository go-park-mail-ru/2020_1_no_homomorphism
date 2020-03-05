package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
)

func MarshallAndWriteProfile(w http.ResponseWriter, profile *models.Profile) {

	profileJson, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		MarshallAndSendError(err, w)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		MarshallAndSendError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CheckAuth(r *http.Request, SessionUC session.UseCase) (*models.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		return nil, err
	}
	user, err := SessionUC.GetUserBySessionID(sid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
