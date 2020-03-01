package delivery

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (api *MyHandler) createCookie(id uuid.UUID) (cookie *http.Cookie) {
	api.Mutex.Lock()
	defer api.Mutex.Unlock()
	sid := uuid.NewV4()
	api.Sessions[sid] = id
	cookie = &http.Cookie{
		Name:     "session_id",
		Value:    sid.String(),
		HttpOnly: true,
		Expires:  time.Now().Add(72 * time.Hour),
	}
	return
}