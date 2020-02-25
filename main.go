package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "no_homomorphism/handlers"
	"no_homomorphism/models"
)

func main() {
	r := mux.NewRouter()
	api := &MyHandler{
		Sessions:     nil,
		UsersStorage: models.NewUsersStorage(),
	}

	fmt.Printf("Starts server at 8080\n")
	r.HandleFunc("/", api.MainHandler)
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", api.LogoutHandler).Methods("DELETE")
	r.HandleFunc("/signup", api.SignUpHandler).Methods("POST")
	r.HandleFunc("/profile/settings", api.SetingsHandler).Methods("PUT")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
