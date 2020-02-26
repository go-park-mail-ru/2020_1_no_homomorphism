package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"no_homomorphism/models"
	"os"
	"strconv"
	"sync"
	"time"
)

type MyHandler struct {
	Sessions     map[uuid.UUID]uuid.UUID // SID -> ID
	UsersStorage *models.UsersStorage
	Mutex        *sync.Mutex
}

func saveFile(file multipart.File, userId string) error {
	fileBody, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to read file body file")
	}
	filePath := "./images/" + userId + ".png" //todo подставлять формат файла
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
		return uuid.FromStringOrNil(""), errors.New(string(http.StatusUnauthorized))
	}
	sessionId, err := uuid.FromString(cookie.Value)
	if err != nil {
		return uuid.FromStringOrNil(""), errors.New(string(http.StatusBadRequest))
	}
	api.Mutex.Lock()
	defer api.Mutex.Unlock()
	userId, ok := api.Sessions[sessionId]
	if !ok {
		return uuid.FromStringOrNil(""), errors.New(string(http.StatusUnauthorized))
	}
	return userId, nil
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
	//
	//mimeType := handler.Header.Get("Content-Type")
	//switch mimeType {
	//case "image/jpeg":
	//case "image/png":
	//default:
	//}
	err = saveFile(file, userId.String())
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
		statusCode, _ := strconv.Atoi(err.Error())
		w.WriteHeader(statusCode)
		fmt.Println(err)
		return
	}

	file, err := os.Open("./images/" + userId.String() + ".png")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	FileHeader := make([]byte, 512)
	_, err = file.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := file.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	w.Header().Set("Content-Disposition", "attachment; filename=profileImage")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	_, _ = file.Seek(0, 0)
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
}

func (api *MyHandler) MainHandler(w http.ResponseWriter, r *http.Request) { //без мьютекса!! - just test
	defer r.Body.Close()
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		id, err := uuid.FromString(session.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, authorized = api.Sessions[id]
	}
	if authorized {
		w.Write([]byte("autrorized"))
	} else {
		w.Write([]byte("not autrorized"))
	}
}

func (api *MyHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(models.UserInput)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userModel, err := api.UsersStorage.GetFullUserInfo(user.Login)
	if err != nil || userModel.Password != user.Password {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Sending status 400 to " + r.RemoteAddr)
		return
	}

	http.SetCookie(w, api.createCookie(userModel.Id))
	w.WriteHeader(http.StatusOK)
	fmt.Println("Sending status 200 to " + r.RemoteAddr)
}

func (api *MyHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sessionId, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || sessionId == nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken, err := uuid.FromString(sessionId.Value)
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
	delete(api.Sessions, sessionToken)
	api.Mutex.Unlock()
	fmt.Println("Mutex Unlocked")

	sessionId.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionId)
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

func (api *MyHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) { //todo чекать на наличие такой записи
	defer r.Body.Close()

	userInput := new(models.UserInput) //todo мб расширить эту структуру
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
	http.SetCookie(w, api.createCookie(userId))
	w.WriteHeader(http.StatusOK)
}

func (api *MyHandler) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
	if newUserData.Password != user.Password {
		log.Print("wrong old password")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	api.UsersStorage.EditUser(user, newUserData)
}
