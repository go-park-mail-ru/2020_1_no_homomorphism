package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/models"
)

var mu = &sync.Mutex{}

var api = MyHandler{
	Sessions: make(map[uuid.UUID]uuid.UUID, 10),
	Mutex:    mu,
	UsersStorage: &models.UsersStorage{
		Users: map[string]*models.User{
			"test": {
				Id:       uuid.NewV4(),
				Login:    "test",
				Password: "$2a$04$0GzSltexrV9gQjFwv5BYuebu7/F13cX.NOupseJQUwqHWDucyBBgO", //123
			},
			"test2": {
				Id:       uuid.NewV4(),
				Login:    "test2",
				Password: "$2a$04$r/rWIhO8ptZAxheWs9cXmeG8fKhICfA5Gko3Qr61ae0.71CwjyODC", //456
			},
			"test3": {
				Id:       uuid.NewV4(),
				Login:    "test3",
				Password: "$2a$04$8G8SC41DvtOYD04qVizzbek.uL9zEI5zlQ3q2Cg.DYekuzMWFsoLa", //789
			},
		},
		Mutex: mu,
	},
}

func TestHandlers_LoginHandler(t *testing.T) {
	jsonUser := bytes.NewBuffer([]byte("{ \"Logindgfd\":\"test\", \"Password\":\"123\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	jsonUser = bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	id, err := uuid.FromString(rr.Result().Cookies()[0].Value)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, api.UsersStorage.Users["test"].Id, api.Sessions[id])
}
func TestHandlers_SignUpHandler(t *testing.T) {
	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test3\", \"Password\":\"111\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.SignUpHandler).ServeHTTP(rr, req)
	assert.Equal(t, rr.Result().StatusCode, http.StatusBadRequest)
	jsonUser = bytes.NewBuffer([]byte("{ \"Logi:\"asdfs\"}"))
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/signup", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.SignUpHandler).ServeHTTP(rr, req)
	assert.Equal(t, rr.Result().StatusCode, http.StatusBadRequest)
	jsonUser = bytes.NewBuffer([]byte("{ \"Login\":\"test4\", \"Password\":\"111\"}"))
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/signup", jsonUser)
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

	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)
	sid := rr.Result().Cookies()[0].Value
	req.Body.Close()
	req, err = http.NewRequest("DELETE", "/logout", bytes.NewBuffer([]byte("")))
	http.HandlerFunc(api.LogoutHandler).ServeHTTP(rr, req)
	assert.Equal(t, rr.Result().StatusCode, http.StatusOK)
	req.AddCookie(rr.Result().Cookies()[0])
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LogoutHandler).ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	id, err := uuid.FromString(sid)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, uuid.FromStringOrNil("0"), api.Sessions[id])
}

func TestMyHandler_SettingsHandler(t *testing.T) {

	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test3\", \"Password\":\"789\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)
	sid, err := uuid.FromString(rr.Result().Cookies()[0].Value)
	req.Body.Close()
	if err != nil {
		t.Error(err)
	}
	id := api.Sessions[sid]
	jsonSettings := bytes.NewBuffer([]byte("{ \"Login\":\"3test\", \"Password\":\"789\",\"new_password\":\"555\" }"))
	req, err = http.NewRequest("PUT", "/profile/settings", jsonSettings)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(rr.Result().Cookies()[0])
	http.HandlerFunc(api.SettingsHandler).ServeHTTP(rr, req)
	idAfter := api.Sessions[uuid.FromStringOrNil(rr.Result().Cookies()[0].Value)]
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, id.String(), idAfter.String())
	assert.Equal(t, api.UsersStorage.Users["3test"].Id, id)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(api.UsersStorage.Users["3test"].Password), []byte("555")))

}

func TestMyHandler_MainHandler(t *testing.T) {
	jsonUser := bytes.NewBuffer([]byte("{}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.MainHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNonAuthoritativeInfo, rr.Result().StatusCode)
	jsonUser = bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)

	req, err = http.NewRequest("GET", "/", jsonUser)
	if err != nil {
		t.Error(err)
	}
	req.AddCookie(rr.Result().Cookies()[0])
	http.HandlerFunc(api.MainHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func TestMyHandler_GetProfileByCookieHandler(t *testing.T) {
	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)
	req, err = http.NewRequest("POST", "/profile/me", jsonUser)
	req.AddCookie(rr.Result().Cookies()[0])
	http.HandlerFunc(api.GetProfileByCookieHandler).ServeHTTP(rr, req)

}

func TestMyHandler_GetProfileHandlerDoesNotExists(t *testing.T) {
	jsonInput := bytes.NewBuffer([]byte("{}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/profiles/abracadabra", jsonInput)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.GetProfileHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func TestMyHandler_PostImageHandler(t *testing.T) {
	jsonInput := bytes.NewBuffer([]byte("{}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/image", jsonInput)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.PostImageHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)
	req, err = http.NewRequest("POST", "/image", jsonUser)
	req.AddCookie(rr.Result().Cookies()[0])
	rr = httptest.NewRecorder()
	http.HandlerFunc(api.PostImageHandler).ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)

}

func Test_GetUserImageHandler(t *testing.T) {
	jsonUser := bytes.NewBuffer([]byte("{ \"Login\":\"test\", \"Password\":\"123\"}"))
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/image", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.GetUserImageHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/login", jsonUser)
	if err != nil {
		t.Error(err)
	}
	http.HandlerFunc(api.LoginHandler).ServeHTTP(rr, req)
	req.AddCookie(rr.Result().Cookies()[0])
	rr = httptest.NewRecorder()
	http.HandlerFunc(api.GetUserImageHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
}

func Test_GetIdByLogin(t *testing.T) {
	id := api.UsersStorage.GetIdByLogin("test")
	assert.Equal(t, api.UsersStorage.Users["test"].Id, id)
}

//
//func TestNewUsersStorage(t *testing.T) {
//	//mu := &sync.Mutex{}
//	_ = models.NewUsersStorage()
//	//assert.NotNil(t, err)
//	_ = models.NewUsersStorage()
//	//assert.Nil(t, err)
//}

func Test_GetUserPassword(t *testing.T) {
	_, err := api.UsersStorage.GetUserPassword("abracadabra")
	assert.NotNil(t, err)
	_, err = api.UsersStorage.GetUserPassword("test")
	assert.Nil(t, err)
}
