package usecase

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	users "no_homomorphism/internal/pkg/user"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"no_homomorphism/internal/pkg/models"
)

type UserUseCase struct {
	Repository users.Repository
	AvatarDir  string
}

var allowedContentType = []string{
	"image/png",
	"image/jpeg",
}

func (uc *UserUseCase) Create(user *models.User) error {
	ok, err := uc.Repository.CheckIfExists(user.Login, user.Email)
	if ok {
		return errors.New("user with this login or email is already exists") //todo сообщать отдельно о логине или\и почте
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	err = uc.Repository.Create(user, hash)
	return err
}

func (uc *UserUseCase) Update(user *models.User, input *models.UserSettings) error {
	var hash []byte
	if input.NewPassword != "" {
		var err error
		hash, err = bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
		if err != nil {
			return err
		}
	}
	return uc.Repository.Update(user, input, hash)
}

func checkFileContentType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil || err == io.EOF {
		return "", err
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	for _, r := range allowedContentType {
		if contentType == r {
			return strings.Split(contentType, "/")[1], nil
		}
	}
	return "", errors.New("this content type is not allowed")
}

func (uc *UserUseCase) UpdateAvatar(user *models.User, file *multipart.FileHeader) (string, error) {

	fileBody, err := file.Open()
	if err != nil {
		return "", errors.New("failed to read file body file")
	}
	defer fileBody.Close()

	contentType, err := checkFileContentType(fileBody)
	if err != nil {
		return "", err
	}

	fileName := user.Id //todo good names for avatars
	filePath := filepath.Join(os.Getenv("MUSIC_PROJ_DIR"), uc.AvatarDir, fileName+"."+contentType)
	fmt.Println(filePath)
	newFile, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("failed to create file")
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, fileBody)
	if err != nil {
		return "", errors.New("error while writing to file")
	}
	err = uc.Repository.UpdateAvatar(user, filepath.Join(uc.AvatarDir, fileName+"."+contentType))
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (uc *UserUseCase) GetUserByLogin(user string) (*models.User, error) {
	return uc.Repository.GetUserByLogin(user)
}

func (uc *UserUseCase) GetProfileByLogin(login string) (*models.Profile, error) {
	user, err := uc.Repository.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}
	profile := &models.Profile{
		Name:  user.Name,
		Login: user.Login,
		Sex:   user.Sex,
		Image: user.Image,
		Email: user.Email,
	}
	return profile, nil
}

func (uc *UserUseCase) GetProfileByUser(user *models.User) *models.Profile {
	profile := &models.Profile{
		Name:  user.Name,
		Login: user.Login,
		Sex:   user.Sex,
		Image: user.Image,
		Email: user.Email,
	}
	return profile
}

func (uc *UserUseCase) CheckUserPassword(user *models.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("wrong password")
	}
	return nil
}
