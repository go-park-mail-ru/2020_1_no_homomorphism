package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"no_homomorphism/internal/pkg/session"
	users "no_homomorphism/internal/pkg/user"

	"no_homomorphism/pkg/logger"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/models"
	"time"
)

type Handler struct {
	SessionDelivery session.Delivery
	UserUC          users.UseCase
	Log             *logger.MainLogger
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(*models.User)
	h.Log.Debug(user)
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
	if r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "user is already auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := &models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	h.Log.Debug(user)
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
	cookie, err := h.SessionDelivery.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Log.LogWarning(r.Context(), "delivery", "Login", "failed to create session: "+err.Error())
		return
	}

	http.SetCookie(w, cookie)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "user is already auth", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input := &models.UserSignIn{}
	err := json.NewDecoder(r.Body).Decode(&input)
	h.Log.Debug(input)
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
	cookie, err := h.SessionDelivery.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Log.LogWarning(r.Context(), "delivery", "Login", "failed to create session: "+err.Error())
		return
	}

	http.SetCookie(w, cookie)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
	err = h.SessionDelivery.Delete(sid)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't delete session:"+err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
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
	h.marshallAndWriteProfile(w, r.Context(), profile)
}

func (h *Handler) SelfProfile(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(*models.User)
	profile := h.UserUC.GetProfileByUser(user)

	h.marshallAndWriteProfile(w, r.Context(), profile)
}

func (h *Handler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
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

	path, err := h.UserUC.UpdateAvatar(user, handler)
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "UpdateAvatar", "failed to update avatar:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.Info("new file created:", path) //add path
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *Handler) marshallAndWriteProfile(w http.ResponseWriter, ctx context.Context, profile *models.Profile) {
	profileJson, err := json.Marshal(profile)
	if err != nil {
		h.Log.HttpInfo(ctx, "error while marshalling:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		h.Log.LogWarning(ctx, "delivery", "marshallAndWriteProfile", "failed to write result"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(ctx, "OK", http.StatusOK)
}
