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

type Suite struct {
	suite.Suite
	redisServer *miniredis.Miniredis
	session     *SessionManager
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

	s.session = NewRedisSessionManager(redisConn)
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
	sID := uuid.NewV4()

	err := s.session.Create(sID.String(), login, uint64(expire.Seconds()))
	require.NoError(s.T(), err)

	value, err := s.redisServer.Get(sID.String())
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, login)

	//test TTL
	s.redisServer.FastForward(expire)

	_, err = s.redisServer.Get(sID.String())
	require.Equal(s.T(), err, errors.New("ERR no such key"))

	//test on closed connection
	s.redisServer.Close()

	err = s.session.Create(sID.String(), login, uint64(expire))
	require.Error(s.T(), err)
}

func (s *Suite) TestDelete() {
	sID := uuid.NewV4()
	testValue := "test_value"
	require.NoError(s.T(), s.redisServer.Set(sID.String(), testValue))

	value, err := s.redisServer.Get(sID.String())
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, testValue)

	require.NoError(s.T(), s.session.Delete(sID.String()))

	_, err = s.redisServer.Get(sID.String())
	require.Equal(s.T(), err, errors.New("ERR no such key"))

	//test on closed connection
	id := uuid.NewV4()

	s.redisServer.Close()
	require.Error(s.T(), s.session.Delete(id.String()))
}

func (s *Suite) TestGetLoginBySessionID() {
	sID := uuid.NewV4()
	testValue := "test_value"
	require.NoError(s.T(), s.redisServer.Set(sID.String(), testValue))

	val, err := s.session.GetLoginBySessionID(sID.String())
	require.NoError(s.T(), err)
	require.Equal(s.T(), testValue, val)

	//test on closed connection
	s.redisServer.Close()

	_, err = s.session.GetLoginBySessionID(sID.String())
	require.Error(s.T(), err)
}
