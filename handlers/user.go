package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/models"
)

type MyHandler struct {
	Sessions     map[uuid.UUID]uuid.UUID // SID -> ID
	UsersStorage *models.UsersStorage
}

//
// func NewMyHandler() *MyHandler {
// 	return &MyHandler{
// 		Sessions: make(map[uuid.UUID]uuid.UUID, 10),
// 		UsersStorage: &models.UsersStorage{
// 			Users: map[string]*models.User{
// 				"test": {
// 					Id:       uuid.FromStringOrNil("1"),
// 					Login: "test",
// 					Password: "123",
// 				},
// 			},
// 			Mutex: sync.RWMutex{},
// 		},
// 	}
// }

func (api *MyHandler) MainHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	authorized := false
	session, err := r.Cookie("session_id")
	mutex := &sync.Mutex{}
	mutex.Lock()
	if err == nil && session != nil {
		id, err := uuid.FromString(session.Value)
		if err != nil {
			mutex.Unlock()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, authorized = api.Sessions[id]
	}
	if authorized {
		w.Write([]byte("autrorized"))
	} else {
		w.Write([]byte("not autrorized"))
	}
	mutex.Unlock()
}

func (api *MyHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(models.UserInput)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mutex := &sync.Mutex{}
	mutex.Lock()
	userModel, ok := api.UsersStorage.Users[user.Login]

	if !ok || userModel.Password != user.Password {
		mutex.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		return
	}

	http.SetCookie(w, api.createCookie(userModel.Id))
	mutex.Unlock()
	w.WriteHeader(http.StatusOK)
	fmt.Println("Sending status 200 to " + r.RemoteAddr)
}

func (api *MyHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userToken, err := uuid.FromString(id.Value)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := api.Sessions[userToken]; !ok {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(api.Sessions, userToken)

	id.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, id)

}

func (api *MyHandler) createCookie(id uuid.UUID) (cookie *http.Cookie) {
	mutex := sync.RWMutex{}
	mutex.Lock()
	sid := uuid.NewV4()
	api.Sessions[sid] = id
	cookie = &http.Cookie{
		Name:    "session_id",
		Value:   sid.String(),
		Expires: time.Now().Add(10 * time.Hour),
	}
	return
}

func (api *MyHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userInput := new(models.UserInput)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := api.UsersStorage.AddUser(userInput)
	if err != nil {
		log.Printf("error while creating User: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, api.createCookie(userId))
	w.WriteHeader(http.StatusBadRequest)
}

func (api *MyHandler) SetingsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusForbidden)
	}
	newUserData := new(models.UserSettings)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newUserData)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusForbidden)
	}
	id := api.Sessions[sid]
	user, err := api.UsersStorage.GetUserById(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if newUserData.Password != user.Password {
		log.Print("wrong old password")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	api.UsersStorage.EditUser(user, newUserData)
}
