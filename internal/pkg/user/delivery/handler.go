package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/user"
)

type Handler struct {
	SessionUC session.UseCase
	UserUC    user.UseCase
	Log       *logger.MainLogger
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool)  {
		h.Log.HttpInfo(r.Context(), "permission denied:"+err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(*models.User)
	input := &models.UserSettings{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Log.HttpInfo(r.Context(), "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.UserUC.Update(user, input); err != nil {
		h.Log.HttpInfo(r.Context(), "can't update user:"+err.Error(), http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("isAuth").(bool)  {
		h.Log.HttpInfo(r.Context(), "already auth"+err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := &models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.UserUC.Create(user); err != nil {
		h.Log.HttpInfo(r.Context(), "error while creating User:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie, err := h.SessionUC.Create(user)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "error while creating cookie:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("isAuth").(bool)  {
		log.Printf("already auth") //todo
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	input := &models.UserSignIn{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.UserUC.GetUserByLogin(input.Login)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't get user from storage: "+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.UserUC.CheckUserPassword(user, input.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Log.HttpInfo(r.Context(), "Login: wrong password", http.StatusBadRequest)
		return
	}
	cookie, err := h.SessionUC.Create(user)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "error while creating session:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool)  {
		log.Printf("permission denied ")//todo
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(*models.User)
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		h.Log.HttpInfo(r.Context(), "could not find cookie", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't get session id from cookie:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.SessionUC.Delete(sid)
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)//check if has errors
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool)  {
		log.Printf("permission denied ")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	login, e := vars["profile"]
	if e == false {
		h.Log.HttpInfo(r.Context(), "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profile, err := h.UserUC.GetProfileByLogin(login)

	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't find profile:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	marshallAndWriteProfile(w, profile)//todo okay dobavit'
}

func (h *Handler) SelfProfile(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool)  {
		log.Printf("permission denied ")//todo
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(*models.User)

	profile := h.UserUC.GetProfileByUser(user)//todo

	marshallAndWriteProfile(w, profile)//todo okay in the end?
}

func (h *Handler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool)  {
		log.Printf("permission denied ")//todo
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(*models.User)

	file, handler, err := r.FormFile("profile_image")
	if err != nil || handler.Size == 0 {
		h.Log.HttpInfo(r.Context(), "can't read profile_image", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file.Close()

	err = h.UserUC.UpdateAvatar(user, handler)
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "UpdateAvatar", "failed to update avatar:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.Info("new file created:", path)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool)  {
		log.Printf("permission denied ")//todo
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

// func (h *Handler) Debug(w http.ResponseWriter, r *http.Request) {
// 	h.UserUC.PrintUserList()
// 	h.SessionUC.PrintSessionList()
// }

func marshallAndWriteProfile(w http.ResponseWriter, profile *models.Profile) {
	profileJson, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)//todo
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)//todo
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}