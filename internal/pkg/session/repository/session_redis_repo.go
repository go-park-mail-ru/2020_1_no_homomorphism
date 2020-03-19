package repository

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

type SessionManager struct {
	redisConn redis.Conn
}

func NewRedisSessionManager(conn redis.Conn) *SessionManager {
	return &SessionManager{
		redisConn: conn,
	}
}

func (sr *SessionManager) Create(sId string, login string, expire time.Duration) error {
	result, err := redis.String(sr.redisConn.Do("SET", sId, login, "EX", int(expire.Seconds())))
	if err != nil {
		return errors.New("failed to write key: " + err.Error())
	}
	if result != "OK" {
		return errors.New("result not OK")
	}
	return nil
}

func (sr *SessionManager) Delete(sId string) error {
	_, err := redis.Int(sr.redisConn.Do("DEL", sId))
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionManager) GetLoginBySessionID(sId string) (string, error) {
	data, err := redis.String(sr.redisConn.Do("GET", sId))
	if err != nil {
		return "", errors.New("cant get data: " + err.Error())
	}
	return data, nil
}
