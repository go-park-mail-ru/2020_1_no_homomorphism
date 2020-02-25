package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"sync"
	"time"
)

//User - model of user
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MyHandler struct {
	Sessions map[string]uuid.UUID
	Users    map[string]*User
	mu       *sync.Mutex
}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		Sessions: make(map[string]uuid.UUID, 10),
		Users: map[string]*User{
			"test": {uuid.FromStringOrNil("1"), "test", "123"},
		},
		mu: &sync.Mutex{},
	}
}

func (api *MyHandler) handler(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	api.mu.Lock()
	if err == nil && session != nil {
		//id, err := uuid.FromString(session.Value)
		//if err != nil {
		//	api.mu.Unlock()
		//	w.WriteHeader(http.StatusBadRequest)
		//	return
		//}
		_, authorized = api.Sessions[session.Value]
	}

	if authorized {
		w.Write([]byte("autrorized"))
	} else {
		w.Write([]byte("not autrorized"))
	}

	api.mu.Unlock()

}

func (api *MyHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(UserInput)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.Lock()
	userModel, ok := api.Users[user.Username]

	if !ok || userModel.Password != user.Password {
		api.mu.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		return
	}

	id := uuid.NewV4()
	api.Sessions[userModel.Username] = id
	api.mu.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   id.String(),
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(200)
	fmt.Println("Sending status 200 to " + r.RemoteAddr)
}

func (api *MyHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, ok := api.Sessions[id.Value]; !ok {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	delete(api.Sessions, id.Value)

	id.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, id)

}

func (api *MyHandler) signUpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userInput := new(UserInput)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInput)

	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.mu.Lock()
	user := User{
		ID:       uuid.NewV1(),
		Username: userInput.Username,
		Password: userInput.Password,
	}
	api.Users[userInput.Username] = &user

	id := uuid.NewV4()
	api.Sessions[id.String()] = id
	api.mu.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   id.String(),
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(200)
}

func main() {
	r := mux.NewRouter()
	api := NewMyHandler()

	fmt.Printf("Starts server at 8080\n")
	r.HandleFunc("/", api.handler)
	r.HandleFunc("/login", api.loginHandler).Methods("POST")
	r.HandleFunc("/logout", api.loginHandler).Methods("DELETE")
	r.HandleFunc("/signup", api.signUpHandler).Methods("POST")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
