package repository

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

type SessionManager struct {
	redisPool redis.Pool
}

func NewRedisSessionManager(conn *redis.Pool) *SessionManager {
	return &SessionManager{
		redisPool: *conn,
	}
}

func (sr *SessionManager) Create(sId string, login string, expire time.Duration) error {
	conn := sr.redisPool.Get()
	result, err := redis.String(conn.Do("SET", sId, login, "EX", int(expire.Seconds())))
	if err != nil {
		return errors.New("failed to write key: " + err.Error())
	}
	if result != "OK" {
		return errors.New("result not OK")
	}
	return nil
}

func (sr *SessionManager) Delete(sId string) error {
	conn := sr.redisPool.Get()
	_, err := redis.Int(conn.Do("DEL", sId))
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionManager) GetLoginBySessionID(sId string) (string, error) {
	conn := sr.redisPool.Get()
	data, err := redis.String(conn.Do("GET", sId))
	if err != nil {
		return "", errors.New("cant get data: " + err.Error())
	}
	return data, nil
}
