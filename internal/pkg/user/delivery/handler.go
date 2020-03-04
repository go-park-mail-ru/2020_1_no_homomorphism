package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/user"
)

type Handler struct {
	SessionUC session.UseCase
	UserUC    user.UseCase
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
	}
	input := &models.UserSettings{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		log.Println("user and session don't match :", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := h.UserUC.Update(user, input); err != nil {
		log.Println("can't update user :", err)
		w.WriteHeader(http.StatusForbidden)
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
		return
	}
	if err := h.UserUC.Create(user); err != nil {
		log.Printf("error while creating User: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, h.SessionUC.Create(user))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	input := &models.UserSignIn{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.UserUC.GetUserByLogin(input.Login)
	if err != nil {
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		log.Println("can't get user from base : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("wrong password : ", err)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		return
	}
	http.SetCookie(w, h.SessionUC.Create(user))
	w.WriteHeader(http.StatusOK)
	fmt.Println("Sending status 200 to " + r.RemoteAddr)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		log.Println("could not find cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Println("can't get session id from cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := h.SessionUC.GetUserBySessionID(sid); err != nil {
		log.Println("this session does not exists : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.SessionUC.Delete(sid)
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Debug(w http.ResponseWriter, r *http.Request) {
	h.UserUC.PrintUserList()
	h.SessionUC.PrintSessionList()
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login, e := vars["profile"]
	if e == false {
		log.Println("no id in mux vars")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profile, err := h.UserUC.GetProfileByLogin(login)

	if err != nil {
		log.Println("can't find this profile :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileJson, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SelfProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		log.Println("could not find cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Println("can't get session id from cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		log.Println("this session does not exists : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profile := h.UserUC.GetProfileByUser(user)

	profileJson, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Header)
	var kek []byte
	_, err := r.Body.Read(kek)
	if err != nil {
		fmt.Println("noooooooooooooooooooo")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		log.Println("could not find cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Println("can't get session id from cookie :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.SessionUC.GetUserBySessionID(sid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		return
	}
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("can't Parse Multipart Form " , err)
		return
	}

	file, handler, err := r.FormFile("profile_image")
	if err != nil || handler.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("can't read profile_image : ", err)
		return
	}
	defer file.Close()

	err = h.UserUC.UpdateAvatar(user, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}


/*

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func saveFile(file multipart.File, userId string, avatarDir string) error {
	fileBody, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to read file body file")
	}
	filePath := os.Getenv("MUSIC_PROJ_DIR") + avatarDir + userId + ".png" // todo подставлять формат файла
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to create file")
	}
	defer newFile.Close()

	_, err = newFile.Write(fileBody)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to write to file")
	}
	return nil
}

func (api *MyHandler) getAvatarPath(userId uuid.UUID) (string, error) {
	// userId, err := api.getUserIdByCookie(r)
	// if err != nil {
	//	return "", err
	// }

	path := api.AvatarDir + userId.String() + ".png"

	isExists, err := exists(os.Getenv("MUSIC_PROJ_DIR") + path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if !isExists {
		log.Println(err)
		return "", errors.New("path does not exists")
	}
	return path, nil
}

func (api *MyHandler) GetTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, e := mux.Vars(r)["id"]
	if e == false {
		log.Println("no id in mux vars")
	}
	requestedID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	track, err := api.TrackStorage.GetFullTrackInfo(uint(requestedID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	writer := json.NewEncoder(w)
	err = writer.Encode(track)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *MyHandler) GetUserImageHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := api.getUserIdByCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	file, err := os.Open(os.Getenv("MUSIC_PROJ_DIR") + api.AvatarDir + userId.String() + ".png")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	FileHeader := make([]byte, 512)
	_, err = file.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	FileSize := strconv.FormatInt(FileStat.Size(), 10) // Get file size as a string

	w.Header().Set("Content-Disposition", "attachment; filename=profileImage")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	_, err = file.Seek(0, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *MyHandler) MainHandler(w http.ResponseWriter, r *http.Request) { // без мьютекса!! - just test
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		id, err := uuid.FromString(session.Value)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, authorized = api.Sessions[id]
	}
	if authorized {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("autrorized"))
	} else {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		w.Write([]byte("not autrorized"))
	}
}

func (api *MyHandler) CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println(http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("SUPER KEEEEEEEEEEEEEK")
		return
	}

	userID, err := uuid.FromString(cookie.Value)
	if err != nil {
		fmt.Println(cookie.Value + "2")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, ok := api.Sessions[userID]; ok {
		fmt.Println(cookie.Value + "3")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Println(cookie.Value + "4")

	return
}


*/
