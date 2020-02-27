package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/models"
)

type MyHandler struct {
	Sessions     map[uuid.UUID]uuid.UUID // SID -> ID
	UsersStorage *models.UsersStorage
	TrackStorage *models.TrackStorage
	Mutex        *sync.Mutex
	AvatarDir    string
}

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
	filePath := avatarDir + userId + ".png" //todo подставлять формат файла
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

func (api *MyHandler) getUserIdByCookie(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		return uuid.FromStringOrNil(""), err
	}
	sessionId, err := uuid.FromString(cookie.Value)
	if err != nil {
		return uuid.FromStringOrNil(""), err
	}
	api.Mutex.Lock()
	defer api.Mutex.Unlock()
	userId, ok := api.Sessions[sessionId]
	if !ok {
		return uuid.FromStringOrNil(""), err
	}
	return userId, nil
}

func (api *MyHandler) getAvatarPath(r *http.Request) (string, error) {
	userId, err := api.getUserIdByCookie(r)
	if err != nil {
		return "", err
	}

	path := api.AvatarDir + userId.String() + ".png"

	isExists, err := exists(path)
	if err != nil {
		return "", err
	}
	if !isExists {
		return "", err
	}
	return path, nil
}

func (api *MyHandler) GetTrackHandler(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.Atoi(mux.Vars(r)["id"])
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

func (api *MyHandler) PostImageHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := api.getUserIdByCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	file, handler, err := r.FormFile("profile_image")
	if err != nil || handler.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = saveFile(file, userId.String(), api.AvatarDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	file, err := os.Open(api.AvatarDir + userId.String() + ".png")
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
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

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

func (api *MyHandler) MainHandler(w http.ResponseWriter, r *http.Request) { //без мьютекса!! - just test
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

func (api *MyHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userModel, err := api.UsersStorage.GetFullUserInfo(user.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(user.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		return
	}

	http.SetCookie(w, api.createCookie(userModel.Id))
	w.WriteHeader(http.StatusOK)
	fmt.Println("Sending status 200 to " + r.RemoteAddr)
}

func (api *MyHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sid, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || sid == nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken, err := uuid.FromString(sid.Value)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.Mutex.Lock()
	if _, ok := api.Sessions[sessionToken]; !ok {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(api.Sessions, uuid.FromStringOrNil(sid.Value))
	api.Mutex.Unlock()
	fmt.Println("Mutex Unlocked")

	sid.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sid)
	w.WriteHeader(http.StatusOK)
}

func (api *MyHandler) createCookie(id uuid.UUID) (cookie *http.Cookie) {
	api.Mutex.Lock()
	defer api.Mutex.Unlock()
	sid := uuid.NewV4()
	api.Sessions[sid] = id
	cookie = &http.Cookie{
		Name:    "session_id",
		Value:   sid.String(),
		Expires: time.Now().Add(10 * time.Hour),
	}
	return
}

func (api *MyHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	userInput := new(models.User)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("New user added")
	http.SetCookie(w, api.createCookie(userId))
	w.WriteHeader(http.StatusOK)
}

func (api *MyHandler) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusForbidden)
	}
	newUserData := new(models.UserSettings)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newUserData)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := uuid.FromString(cookie.Value)
	if err != nil {
		log.Printf("permission denied: %s", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	api.Mutex.Lock()
	id := api.Sessions[sid]
	api.Mutex.Unlock()
	user, err := api.UsersStorage.GetUserById(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUserData.Password)); err != nil {
		log.Print("wrong old password")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = api.UsersStorage.EditUser(user, newUserData)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *MyHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["profile"]

	fmt.Println(login)

	profile, err := api.UsersStorage.GetProfileByLogin(login)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path, err := api.getAvatarPath(r)
	if err != nil {
		fmt.Println(err)
		profile.Image = api.AvatarDir + "default.png"
	} else {
		profile.Image = path
	}
	profile.Image = path
	profileJson, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (api *MyHandler) GetProfileByCookieHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := api.getUserIdByCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	user, err := api.UsersStorage.GetUserById(userId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	profile, err := api.UsersStorage.GetProfileByLogin(user.Login)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path, err := api.getAvatarPath(r)
	if err != nil {
		fmt.Println(err)
		profile.Image = api.AvatarDir + "default.png"
	} else {
		profile.Image = path
	}
	profileJson, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(profileJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
