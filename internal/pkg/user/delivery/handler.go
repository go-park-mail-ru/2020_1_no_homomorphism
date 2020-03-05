package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/middleware"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/user"
)

type Handler struct {
	SessionUC session.UseCase
	UserUC    user.UseCase
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	user, err := middleware.CheckAuth(r, h.SessionUC)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		middleware.MarshallAndSendError(err, w)
		return
	}

	input := &models.UserSettings{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}

	if err := h.UserUC.Update(user, input); err != nil {
		log.Println("can't update user :", err)
		w.WriteHeader(http.StatusForbidden)
		middleware.MarshallAndSendError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	if err := h.UserUC.Create(user); err != nil {
		log.Printf("error while creating User: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	cookie, err := h.SessionUC.Create(user)
	if err != nil {
		log.Printf("error while creating cookie: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		middleware.MarshallAndSendError(err, w)
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &models.UserSignIn{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	user, err := h.UserUC.GetUserByLogin(input.Login)
	if err != nil {
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		log.Println("can't get user from base : ", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	if err := h.UserUC.CheckUserPassword(user, input.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		log.Println("wrong password : ", err)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
	}
	cookie, err := h.SessionUC.Create(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		middleware.MarshallAndSendError(err, w)
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	fmt.Println("Sending status 200 to " + r.RemoteAddr)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		log.Println("could not find cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Println("can't get session id from cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	h.SessionUC.Delete(sid)
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.CheckAuth(r, h.SessionUC)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		middleware.MarshallAndSendError(err, w)
		return
	}
	vars := mux.Vars(r)
	login, e := vars["profile"]
	if e == false {
		log.Println("no id in mux vars")
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	profile, err := h.UserUC.GetProfileByLogin(login)

	if err != nil {
		log.Println("can't find this profile :", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	middleware.MarshallAndWriteProfile(w, profile)
}

func (h *Handler) SelfProfile(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.CheckAuth(r, h.SessionUC)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		middleware.MarshallAndSendError(err, w)
		return
	}
	if err != nil {
		log.Println("this session does not exists : ", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	profile := h.UserUC.GetProfileByUser(user)

	middleware.MarshallAndWriteProfile(w, profile)
}

func (h *Handler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {

	_, err := middleware.CheckAuth(r, h.SessionUC)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		middleware.MarshallAndSendError(err, w)
		return
	}
	var buffer []byte
	_, err = r.Body.Read(buffer)
	if err != nil {
		log.Println("can't read request body")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		log.Println("could not find cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Println("can't get session id from cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		middleware.MarshallAndSendError(err, w)
		return
	}
	user, err := h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		middleware.MarshallAndSendError(err, w)
		return
	}
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("can't Parse Multipart Form ", err)

		middleware.MarshallAndSendError(err, w)
		return
	}

	file, handler, err := r.FormFile("profile_image")
	if err != nil || handler.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("can't read profile_image : ", err)
		middleware.MarshallAndSendError(err, w)
		return
	}
	defer file.Close()

	err = h.UserUC.UpdateAvatar(user, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		middleware.MarshallAndSendError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.CheckAuth(r, h.SessionUC)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		middleware.MarshallAndSendError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// func (h *Handler) Debug(w http.ResponseWriter, r *http.Request) {
// 	h.UserUC.PrintUserList()
// 	h.SessionUC.PrintSessionList()
// }
