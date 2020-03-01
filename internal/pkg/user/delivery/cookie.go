package delivery

import (
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func  CreateCookie(sid uuid.UUID)  *http.Cookie {
	return  &http.Cookie{
		Name:     "session_id",
		Value:    sid.String(),
		HttpOnly: true,
		Expires:  time.Now().Add(72 * time.Hour),
	}
}