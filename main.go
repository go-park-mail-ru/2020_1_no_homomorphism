package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "no_homomorphism/handlers"
	"no_homomorphism/models"
)

type MyHandler struct {
	Sessions     map[uuid.UUID]uuid.UUID // SID -> ID
	UsersStorage models.UsersStorage
	Mutex sync.Mutex
}

//func (api *MyHandler) getUserIdByCookie(r *http.Request) (uuid.UUID, error) {
//	session, err := r.Cookie("session_id")
//	if err == http.ErrNoCookie {
//		return uuid.FromStringOrNil(""), errors.New(string(http.StatusUnauthorized))
//	}
//	sessionId, err := uuid.FromString(session.Value)
//	if err != nil {
//		return uuid.FromStringOrNil(""), errors.New(string(http.StatusBadRequest))
//	}
//
//	userId := api.Sessions[sessionId]
//}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		Sessions: make(map[uuid.UUID]uuid.UUID, 10),
		UsersStorage: models.UsersStorage{
			Users: map[string]*models.User{
				"test": {
					Id:       uuid.FromStringOrNil("1"),
					Nickname: "test",
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mutex := &sync.Mutex{}
	mutex.Lock()
	userModel, ok := api.UsersStorage.Users[user.Nickname]

	if !ok || userModel.Password != user.Password {
		mutex.Unlock()
		w.WriteHeader(400)
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userToken, err := uuid.FromString(id.Value)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, ok := api.Sessions[userToken]; !ok {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
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
		w.WriteHeader(http.StatusBadRequest)
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

func (api *MyHandler) editUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(403)
	}
	newUserData := new(models.User)
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
	api.UsersStorage.EditUser(user, newUserData)
}

// func (api *MyHandler) showHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	for r, w := range api.UsersStorage.Users {
// 		fmt.Println(r, w)
// 	}
// }
func (api *MyHandler) saveImageHandler(w http.ResponseWriter, r *http.Request) {
	//userId := api.getUserIdByCookie()
	//
	//err = r.ParseMultipartForm(5 * 1024 * 1025)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	fmt.Println(err)
	//	return
	//}

	file, handler, err := r.FormFile("my_file")
	if err != nil || handler.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileBody, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	filePath := "./images/" + uuid.NewV4().String() + ".png" //todo подставлять формат файла
	newFile, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	_, err = newFile.Write(fileBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
