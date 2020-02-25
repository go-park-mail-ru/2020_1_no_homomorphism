package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"no_homomorphism/models"
)


func TestHandlers_LoginHandler(t *testing.T) {
	 api := MyHandler{Sessions: make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: &models.UsersStorage{
			Users: map[string]*models.User{
				"test": {
					Id:       uuid.FromStringOrNil("1"),
					Login:    "test",
					Password: "123",
				},
				"test2": {
					Id:       uuid.FromStringOrNil("2"),
					Login:    "test2",
					Password: "456",
				},
				"test3": {
					Id:       uuid.FromStringOrNil("3"),
					Login:    "test3",
					Password: "789",
				},
			},
			Mutex: sync.RWMutex{},
		},
	}
	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	id, err := uuid.FromString(rr.Result().Cookies()[0].Value)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, api.UsersStorage.Users["test"].Id, api.Sessions[id])
}
func TestHandlers_SignUpHandler(t *testing.T) {
	api := MyHandler{Sessions: make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: &models.UsersStorage{
			Users: map[string]*models.User{
				"test": {
					Id:       uuid.FromStringOrNil("1"),
					Login:    "test",
					Password: "123",
				},
				"test2": {
					Id:       uuid.FromStringOrNil("2"),
					Login:    "test2",
					Password: "456",
				},
				"test3": {
					Id:       uuid.FromStringOrNil("3"),
					Login:    "test3",
					Password: "789",
				},
			},
			Mutex: sync.RWMutex{},
		},
	}
	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test4\", \"Password\":\"111\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.SignUpHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	id, err := uuid.FromString(rr.Result().Cookies()[0].Value)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, api.UsersStorage.Users["test4"].Id, api.Sessions[id])
}

func TestHandlers_LogoutHandler(t *testing.T) {
	api := MyHandler{Sessions: make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: &models.UsersStorage{
			Users: map[string]*models.User{
				"test": {
					Id:       uuid.FromStringOrNil("1"),
					Login:    "test",
					Password: "123",
				},
				"test2": {
					Id:       uuid.FromStringOrNil("2"),
					Login:    "test2",
					Password: "456",
				},
				"test3": {
					Id:       uuid.FromStringOrNil("3"),
					Login:    "test3",
					Password: "789",
				},
			},
			Mutex: sync.RWMutex{},
		},
	}
	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)
	fmt.Println(rr.Result().Cookies()[0].Value)
	req, err = http.NewRequest("DELETE", "/logout", nil)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LogoutHandler).ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

	assert.Nil(t, rr.Result().Cookies()[0].Value)
}