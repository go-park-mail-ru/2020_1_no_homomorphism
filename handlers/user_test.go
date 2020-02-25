package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify"
	"no_homomorphism/models"
)


func TestHandlers_LoginHandler(t *testing.T) {
	 api := &MyHandler{Sessions: make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: &models.UsersStorage{
			Users: map[string]*models.User{
				"test": {
					Id:       uuid.FromStringOrNil("1"),
					Login:    "test",
					Password: "123",
				},
				"test2": {
					Id:       uuid.FromStringOrNil("1"),
					Login:    "test2",
					Password: "456",
				},
				"test3": {
					Id:       uuid.FromStringOrNil("1"),
					Login:    "test2",
					Password: "789",
				},
			},
			Mutex: sync.RWMutex{},
		},
	}
	jsonUser := bytes.NewBuffer([]byte("{ 'Login':'test2', 'Password':'456'}"))
	_, err := http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	//loginHandler := http.HandleFunc(api.LoginHandler)





}