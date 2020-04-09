package session

import (
	"time"
)

type Repository interface {
	Create(sId string, value string, expire time.Duration) error
	Delete(sId string) error
	GetLoginBySessionID(sId string) (string, error)
}
