package repository

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

type SessionManager struct {
	redisConn redis.Conn
}

func NewRedisSessionManager(conn redis.Conn) *SessionManager {
	return &SessionManager{
		redisConn: conn,
	}
}

func (sr *SessionManager) Create(login string, expire time.Duration) (uuid.UUID, error) {
	newUUID := uuid.NewV4()//todo replace uuid gen
	mKey := "sessions:" + newUUID.String()
	result, err := redis.String(sr.redisConn.Do("SET", mKey, login, "EX", int(expire.Seconds())))
	if err != nil {
		return uuid.UUID{}, errors.New("failed to write key: " + err.Error())
	}
	if result != "OK" {
		return uuid.UUID{}, errors.New("result not OK")
	}
	return newUUID, nil
}

func (sr *SessionManager) Delete(sessionID uuid.UUID) error {
	mKey := "sessions:" + sessionID.String()
	_, err := redis.Int(sr.redisConn.Do("DEL", mKey))
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionManager) GetLoginBySessionID(sessionID uuid.UUID) (string, error) {
	mKey := "sessions:" + sessionID.String()
	data, err := redis.String(sr.redisConn.Do("GET", mKey))
	if err != nil {
		return "", errors.New("cant get data: " + err.Error())
	}
	return data, nil
}
