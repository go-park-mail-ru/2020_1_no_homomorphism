package delivery

import (
	"context"
	"encoding/json"
	"net/http"

	"no_homomorphism/internal/pkg/csrf"
	"no_homomorphism/internal/pkg/session"
	users "no_homomorphism/internal/pkg/user"

	"no_homomorphism/pkg/logger"

	"time"

	"github.com/gorilla/mux"
	"no_homomorphism/internal/pkg/models"
)

type UserHandler struct {
	SessionDelivery session.Delivery
	UserUC          users.UseCase
	CSRF            csrf.CryptToken
	Log             *logger.MainLogger
	ImgTypes        map[string]string
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !r.Context().Value("isCSRFTokenCorrect").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user has wrong csrf token", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := r.Context().Value("user").(models.User)

	input := models.UserSettings{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Log.HttpInfo(r.Context(), "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	emailExists, err := h.UserUC.Update(user, input)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't update user:" + err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if emailExists != users.NO {
		h.Log.HttpInfo(r.Context(), "user with same email exists", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "user is already auth", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exists, err := h.UserUC.Create(user)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "error while creating User:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if exists != users.NO {
		h.checkAndSendExisting(w, r.Context(), exists)
		return
	}
	cookie, err := h.SessionDelivery.Create(user)
	if err != nil {
		h.Log.LogWarning(r.Context(), "user delivery", "Create", "failed to create session: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := h.CSRF.Create(cookie.Value, cookie.Expires.Unix())
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "Login", "failed to create csrf token: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "Csrf-Token")
	w.Header().Set("Csrf-Token", token)

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusCreated)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusCreated)
}

func (h *UserHandler) checkAndSendExisting(w http.ResponseWriter, ctx context.Context, exists users.SameUserExists) {
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	type createResponse struct {
		LoginExists bool `json:"login_exists"`
		EmailExists bool `json:"email_exists"`
	}
	response := createResponse{false, false}
	switch exists {
	case users.EMAIL:
		response.EmailExists = true
	case users.LOGIN:
		response.LoginExists = true
	case users.FULL:
		response.EmailExists = true
		response.LoginExists = true
	}
	h.Log.HttpInfo(ctx, "user with same data is already exists", http.StatusConflict)
	w.WriteHeader(http.StatusConflict)

	err := writer.Encode(response)
	if err != nil {
		h.Log.LogWarning(ctx, "user delivery", "checkAndSendExisting", "failed to encode: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "user is already auth", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	input := models.UserSignIn{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "error while unmarshalling JSON:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userData, err := h.UserUC.Login(input)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to login:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie, err := h.SessionDelivery.Create(userData)
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "Login", "failed to create session: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := h.CSRF.Create(cookie.Value, cookie.Expires.Unix())
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "Login", "failed to create csrf token: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "Csrf-Token")
	w.Header().Set("Csrf-Token", token)

	http.SetCookie(w, &cookie)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !r.Context().Value("isCSRFTokenCorrect").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user has wrong csrf token", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		h.Log.HttpInfo(r.Context(), "could not find cookie", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.SessionDelivery.Delete(cookie.Value)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "can't delete session:"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) { //todo достпуность данных других пользователей??
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
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(&profile)
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "selfProfile", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *UserHandler) SelfProfile(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(models.User)
	profile := h.UserUC.GetOutputUserData(user)

	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err := writer.Encode(&profile)
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "selfProfile", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !r.Context().Value("isCSRFTokenCorrect").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user has wrong csrf token", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := r.Context().Value("user").(models.User) //todo check error

	file, handler, err := r.FormFile("profile_image")
	if err != nil || handler.Size == 0 {
		h.Log.HttpInfo(r.Context(), "can't read profile_image", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := handler.Header.Get("Content-Type")
	elem, ok := h.ImgTypes[mimeType]
	if !ok {
		h.Log.HttpInfo(r.Context(), "wrong file content-type", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path, err := h.UserUC.UpdateAvatar(user, file, elem)
	if err != nil {
		h.Log.LogWarning(r.Context(), "delivery", "UpdateAvatar", "failed to update avatar:"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.Info("new file created:", path)
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *UserHandler) GetUserStat(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		h.Log.HttpInfo(r.Context(), "no data in mux vars", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userStat, err := h.UserUC.GetUserStat(id)
	if err != nil {
		h.Log.HttpInfo(r.Context(), "failed to get user's stat"+err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(userStat)
	if err != nil {
		h.Log.LogWarning(r.Context(), "user delivery", "GetUserStat", "failed to encode json"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}

func (h *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value("isAuth").(bool) {
		h.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.Log.HttpInfo(r.Context(), "OK", http.StatusOK)
}
