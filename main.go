package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/models"
)

type MyHandler struct {
	Sessions     map[uuid.UUID]uuid.UUID // SID -> ID
	UsersStorage models.UsersStorage
}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		Sessions: make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: models.UsersStorage{
			Users: map[string]*models.User{
				"test": {
					Id:       uuid.FromStringOrNil("1"),
					Login: "test",
					Password: "123",
				},
			},
			Mutex: sync.RWMutex{},
		},
	}
}

func (api *MyHandler) handler(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	mutex := &sync.Mutex{}
	mutex.Lock()
	if err == nil && session != nil {
		id, err := uuid.FromString(session.Value)
		if err != nil {
			mutex.Unlock()
			w.WriteHeader(400)
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

func (api *MyHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(models.UserInput)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(400)
		return
	}

	mutex := &sync.Mutex{}
	mutex.Lock()
	userModel, ok := api.UsersStorage.Users[user.Login]

	if !ok || userModel.Password != user.Password {
		mutex.Unlock()
		w.WriteHeader(400)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		return
	}

	http.SetCookie(w, api.createCookie(userModel.Id))
	mutex.Unlock()
	w.WriteHeader(200)
	fmt.Println("Sending status 200 to " + r.RemoteAddr)
}

func (api *MyHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	userToken, err := uuid.FromString(id.Value)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	if _, ok := api.Sessions[userToken]; !ok {
		fmt.Println(err)
		w.WriteHeader(401)
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

func (api *MyHandler) signUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userInput := new(models.UserInput)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(400)
		return
	}
	userId, err := api.UsersStorage.AddUser(userInput)
	if err != nil {
		log.Printf("error while creating User: %s", err)
		w.WriteHeader(400)
		return
	}
	http.SetCookie(w, api.createCookie(userId))
	w.WriteHeader(200)
}

func (api *MyHandler) setingsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(403)
	}
	newUserData := new(models.UserSettings)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newUserData)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(400)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(403)
	}
	id := api.Sessions[sid]
	user, err := api.UsersStorage.GetById(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(403)
		return
	}
	if newUserData.Password != user.Password {
		log.Print("wrong old password")
		w.WriteHeader(403)
		return
	}
	api.UsersStorage.EditUser(user, newUserData)
}

// func (api *MyHandler) showHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	for r, w := range api.UsersStorage.Users {
// 		fmt.Println(r, w)
// 	}
// }
func main() {
	r := mux.NewRouter()
	api := NewMyHandler()

	fmt.Printf("Starts server at 8080\n")
	r.HandleFunc("/", api.handler)
	r.HandleFunc("/login", api.loginHandler).Methods("POST")
	r.HandleFunc("/logout", api.loginHandler).Methods("DELETE")
	r.HandleFunc("/signup", api.signUpHandler).Methods("POST")
	// r.HandleFunc("/show", api.showHandler).Methods("POST")
	r.HandleFunc("/profile/settings", api.setingsHandler).Methods("POST")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
