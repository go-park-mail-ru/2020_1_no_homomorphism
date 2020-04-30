package repository

import (
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	redisServer *miniredis.Miniredis
	session     TokenManager
	tokenValue  string
}

func (s *Suite) SetupSuite() {
	var err error
	s.redisServer, err = miniredis.Run()
	require.NoError(s.T(), err)

	addr := s.redisServer.Addr()
	redisConn := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
	s.tokenValue = "1"

	s.session = NewRedisTokenManager(redisConn)
}

//Need to restore connection after each func with closed connection testing
func (s *Suite) AfterTest(_, _ string) {
	s.SetupSuite()
}

func (s *Suite) TearDownSuite() {
	s.redisServer.Close()
}

func TestSessions(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestAdd() {
	token := "testToken"
	expire := 3600

	err := s.session.Add(token, int64(expire))
	require.NoError(s.T(), err)

	value, err := s.redisServer.Get(token)
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, s.tokenValue)

	//test TTL
	s.redisServer.FastForward(time.Second * 4000)

	_, err = s.redisServer.Get(token)
	require.Equal(s.T(), err, errors.New("ERR no such key"))

	//test on closed connection
	s.redisServer.Close()

	err = s.session.Add(token, int64(expire))
	require.Error(s.T(), err)
}

func (s *Suite) TestGetLoginBySessionID() {
	token := "testToken"
	require.NoError(s.T(), s.redisServer.Set(token, s.tokenValue))

	err := s.session.Check(token)
	require.Error(s.T(), err)

	//test no token
	newToken := "wer2v2v3"

	err = s.session.Check(newToken)
	require.NoError(s.T(), err)

	//test on closed connection
	s.redisServer.Close()

	err = s.session.Check(token)
	require.Error(s.T(), err)
}
