package repository

import (
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

//go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

type Suite struct {
	suite.Suite
	redisServer *miniredis.Miniredis
	session     *SessionManager
	keyPrefix   string
}

func (s *Suite) SetupSuite() {
	var err error
	s.redisServer, err = miniredis.Run()
	require.NoError(s.T(), err)

	conn, err := redis.Dial("tcp", s.redisServer.Addr())
	s.session = NewRedisSessionManager(conn)
	s.keyPrefix = "sessions:"
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

func (s *Suite) TestCreate() {
	login := "test_login"
	expire := time.Hour * 8
	sId, err := s.session.Create(login, expire)
	require.NoError(s.T(), err)

	value, err := s.redisServer.Get(s.keyPrefix + sId.String())
	require.NoError(s.T(), err)

	require.Equal(s.T(), value, login)

	//test TTL
	s.redisServer.FastForward(expire)

	_, err = s.redisServer.Get(s.keyPrefix + sId.String())
	require.Equal(s.T(), err, errors.New("ERR no such key"))

	//test on closed connection
	s.redisServer.Close()

	_, err = s.session.Create(login, expire)
	require.Error(s.T(), err)
}

func (s *Suite) TestDelete() {
	sId := uuid.NewV4()
	testValue := "test_value"
	require.NoError(s.T(), s.redisServer.Set(s.keyPrefix+sId.String(), testValue))

	value, err := s.redisServer.Get(s.keyPrefix + sId.String())
	require.Equal(s.T(), value, testValue)

	require.NoError(s.T(), s.session.Delete(sId))

	value, err = s.redisServer.Get(s.keyPrefix + sId.String())
	require.Equal(s.T(), err, errors.New("ERR no such key"))

	//test on closed connection
	id := uuid.NewV4()

	s.redisServer.Close()
	require.Error(s.T(), s.session.Delete(id))
}

func (s *Suite) TestGetLoginBySessionID() {
	sId := uuid.NewV4()
	testValue := "test_value"
	require.NoError(s.T(), s.redisServer.Set(s.keyPrefix+sId.String(), testValue))

	val, err := s.session.GetLoginBySessionID(sId)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testValue, val)

	//test on closed connection
	s.redisServer.Close()

	_, err = s.session.GetLoginBySessionID(sId)
	require.Error(s.T(), err)
}
