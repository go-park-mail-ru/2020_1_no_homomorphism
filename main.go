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
	fmt.Println("New request from " + r.RemoteAddr)
	w.Write([]byte("Hello,"))
	fmt.Fprintf(w, r.RemoteAddr+"\n")
}

func (api *MyHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New request from " + r.RemoteAddr)
	defer r.Body.Close()

	user := new(UserInput)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{status: 500}"))
		return
	}

	fmt.Println("i am here")

	api.mu.Lock()
	userModel, ok := api.Users[user.Username]

	if !ok {
		api.mu.Unlock()
		w.WriteHeader(400)
		return
	}

	if userModel.Password != user.Password {
		api.mu.Unlock()
		w.WriteHeader(400)
		return
	}
	fmt.Println("i am here2")

	id := uuid.NewV4()

	api.Sessions[userModel.Username] = id

	api.mu.Unlock()
	fmt.Println("i am here3")

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

	fmt.Printf("Starts server at 8080")
	r.HandleFunc("/", api.handler)
	r.HandleFunc("/login", api.loginHandler).Methods("POST")
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
