package repository

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type TokenManager struct {
	redisPool *redis.Pool
}

func NewRedisTokenManager(conn *redis.Pool) TokenManager {
	return TokenManager{
		redisPool: conn,
	}
}

func (sr *TokenManager) Add(token string, expire int64) error {
	conn := sr.redisPool.Get()
	result, err := redis.String(conn.Do("SET", token, 1, "EX", expire))
	if err != nil {
		return errors.New("failed to write key: " + err.Error())
	}
	if result != "OK" {
		return errors.New("result not OK")
	}
	return nil
}

func (sr *TokenManager) Check(token string) error {
	conn := sr.redisPool.Get()
	_, err := redis.String(conn.Do("GET", token))

	if err != nil {
		if err == redis.ErrNil {
			return nil
		}
		return err
	}
	return fmt.Errorf("token is not valid")
}
