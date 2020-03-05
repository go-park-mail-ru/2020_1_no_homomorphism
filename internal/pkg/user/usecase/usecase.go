package usecase

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/user"
)

var mutex = &sync.Mutex{}

type UserUseCase struct {
	Repository user.Repository
	AvatarDir  string
}

var allowedContentType = []string{
	"image/png",
	"image/jpeg",
	"image/jpg",
}

func (uc *UserUseCase) Create(user *models.User) error {
	_, ok := uc.GetUserByLogin(user.Login)
	if ok == nil {
		return errors.New("user with this login is already exists")
	}
	err := uc.Repository.Create(user)
	uc.Repository.UpdateAvatar(user, uc.AvatarDir+"default.png")
	return err
}

func (uc *UserUseCase) Update(user *models.User, input *models.UserSettings) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return errors.New("old password is wrong")
	}
	return uc.Repository.Update(user, input)
}

func getFileContentType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	contentType := http.DetectContentType(buffer[:n])
	fmt.Println("CONTENT TYPE: ", contentType)
	return contentType, nil
}

func checkFileContentType(file multipart.File) (string, error) {
	contentType, err := getFileContentType(file)
	if err != nil {
		return "", err
	}
	for _, r := range allowedContentType {
		if contentType == r {
			return strings.Split(contentType, "/")[1], nil
		}
	}
	return "", errors.New("this content type does not allowed")
}

func (uc *UserUseCase) UpdateAvatar(user *models.User, file multipart.File) error {
	fmt.Println("HELLO")
	fileBody, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return errors.New("failed to read file body file")
	}
	//contentType, err := checkFileContentType(file)
	contentType := "png"
	if err != nil {
		log.Println("error while checking content type :", err)
		return err
	}
	fileName := strconv.Itoa(int(user.Id)) //todo good names for avatars
	filePath := os.Getenv("MUSIC_PROJ_DIR") + uc.AvatarDir + fileName + "." + contentType
	fmt.Println(filePath)
	newFile, err := os.Create(filePath)
	if err != nil {
		log.Println("failed to create file", err)
		return errors.New("failed to create file")
	}
	defer newFile.Close()
	_, err = newFile.Write(fileBody)
	if err != nil {
		log.Println("error while writing to file", err)
		return errors.New("error while writing to file")
	}
	uc.Repository.UpdateAvatar(user, uc.AvatarDir+fileName+"."+contentType)
	return nil
}

func (uc *UserUseCase) GetUserByLogin(user string) (*models.User, error) {
	return uc.Repository.GetUserByLogin(user)
}

func (uc *UserUseCase) GetProfileByLogin(login string) (*models.User, error) {
	user, err := uc.Repository.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}
	profile := &models.User{
		Name:  user.Name,
		Login: user.Login,
		Sex:   user.Sex,
		Image: user.Image,
		Email: user.Email,
	}
	return profile, nil
}

func (uc *UserUseCase) GetProfileByUser(user *models.User) *models.User {
	profile := &models.User{
		Name:  user.Name,
		Login: user.Login,
		Sex:   user.Sex,
		Image: user.Image,
		Email: user.Email,
	}
	return profile
}

func (uc *UserUseCase) PrintUserList() {
	uc.Repository.PrintUserList()
}
