package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/user"
	"no_homomorphism/pkg/logger"
)

type Handler struct {
	SessionUC session.UseCase
	UserUC    user.UseCase
	Log       *logger.MainLogger
}

//todo поставить id запроса всем логам

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.Log.HttpInfo("", "permission denied:"+err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	input := &models.UserSettings{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Log.HttpInfo("", "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		h.Log.HttpInfo("", "permission denied:"+err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		h.Log.HttpInfo("", "user and session don't match:"+err.Error(), http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := h.UserUC.Update(user, input); err != nil {
		h.Log.HttpInfo("", "can't update user:"+err.Error(), http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	h.Log.HttpInfo("", "OK", http.StatusOK)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		h.Log.HttpInfo("", "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.UserUC.Create(user); err != nil {
		h.Log.HttpInfo("", "error while creating User:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, h.SessionUC.Create(user))
	h.Log.HttpInfo("", "OK", http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.Log.LogRequest(*r)
	input := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.Log.HttpInfo("", "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.UserUC.GetUserByLogin(input.Login)
	if err != nil {
		h.Log.HttpInfo("", "can't get user from storage: "+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		h.Log.HttpInfo("", "Login: wrong password", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, h.SessionUC.Create(user))
	h.Log.HttpInfo("", "OK", http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		h.Log.HttpInfo("", "could not find cookie", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		h.Log.HttpInfo("", "can't get session id from cookie:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := h.SessionUC.GetUserBySessionID(sid); err != nil {
		h.Log.HttpInfo("", "session doesn't exists:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.SessionUC.Delete(sid) //todo handle error
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	h.Log.HttpInfo("", "OK", http.StatusOK)
}

func (h *Handler) Debug(w http.ResponseWriter, r *http.Request) {
	h.UserUC.PrintUserList()
	h.SessionUC.PrintSessionList()
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login, ok := vars["profile"]
	if !ok {
		h.Log.HttpInfo("", "no id in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profile, err := h.UserUC.GetProfileByLogin(login)
	if err != nil {
		h.Log.HttpInfo("", "can't find profile:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileJson, err := json.Marshal(profile)
	if err != nil {
		h.Log.HttpInfo("", "Error while marshalling:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		h.Log.LogWarning("", "delivery", "Profile", "failed to write:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo("", "OK", http.StatusOK)
}

func (h *Handler) SelfProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		h.Log.HttpInfo("", "no cookie", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		h.Log.HttpInfo("", "can't get session id from cookie:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		h.Log.HttpInfo("", "this session does not exists: "+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profile := h.UserUC.GetProfileByUser(user)

	profileJson, err := json.Marshal(profile)
	if err != nil {
		h.Log.HttpInfo("", "Error while marshalling:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		h.Log.LogWarning("", "delivery", "SelfProfile", "failed to write:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo("", "OK", http.StatusOK)
}

func (h *Handler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	var buff []byte
	_, err := r.Body.Read(buff)
	if err != nil {
		h.Log.LogWarning("", "delivery", "UpdateAvatar", "failed to read body:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		h.Log.HttpInfo("", "no cookie", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		h.Log.HttpInfo("", "can't get session id from cookie:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		h.Log.HttpInfo("", "this session does not exists: "+err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.Log.HttpInfo("", "can't Parse Multipart Form:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("profile_image")
	if err != nil || handler.Size == 0 {
		h.Log.HttpInfo("", "can't read profile_image", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.UserUC.UpdateAvatar(user, file)
	if err != nil {
		h.Log.LogWarning("", "delivery", "UpdateAvatar", "failed to update avatar:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo("", "OK", http.StatusOK)
}

func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.Log.HttpInfo("", "no cookie", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		h.Log.HttpInfo("", "can't get session id from cookie:"+err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		h.Log.HttpInfo("", "this session does not exists: "+err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.Log.HttpInfo("", "OK", http.StatusOK)
}
