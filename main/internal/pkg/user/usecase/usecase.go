package usecase

import (
	"context"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/proto/filetransfer"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"path/filepath"

	users "github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/user"
)

type UserUseCase struct {
	Repository  users.Repository
	FileService filetransfer.UploadServiceClient
	AvatarDir   string
}

func (uc *UserUseCase) Create(user models.User) (users.SameUserExists, error) {
	loginExists, emailExists, err := uc.Repository.CheckIfExists(user.Login, user.Email)
	if err != nil {
		return users.FULL, err
	}
	if loginExists && emailExists {
		return users.FULL, nil
	}
	if loginExists {
		return users.LOGIN, nil
	}
	if emailExists {
		return users.EMAIL, nil
	}
	return users.NO, uc.Repository.Create(user)
}

func (uc *UserUseCase) Update(user models.User, input models.UserSettings) (users.SameUserExists, error) {
	if user.Email != input.Email {
		_, emailExists, err := uc.Repository.CheckIfExists("", input.Email)
		if err != nil {
			return users.FULL, fmt.Errorf("failed to check email existing: %v", err)
		}
		if emailExists {
			return users.EMAIL, nil
		}
	}
	return users.NO, uc.Repository.Update(user, input)
}

func (uc *UserUseCase) UpdateAvatar(user models.User, file io.Reader, fileType string) (string, error) {
	fileName := uuid.NewV4().String()
	fullFileName := fileName + "." + fileType

	md := metadata.New(map[string]string{"fileName": filepath.Join(os.Getenv("FILE_ROOT")+uc.AvatarDir, fullFileName)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := uc.FileService.Upload(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	write := true
	chunk := make([]byte, 1024)

	for write {
		size, err := file.Read(chunk)
		if err != nil {
			if err == io.EOF {
				write = false
				continue
			}
			return "", fmt.Errorf("failed to read file: %v", err)
		}
		err = stream.Send(&filetransfer.Chunk{Content: chunk[:size]})
		if err != nil {
			return "", fmt.Errorf("failed to send file to service: %v", err)
		}
	}

	status, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("error occured in filetransfer service: %v, status: %v", err, status)
	}

	return uc.Repository.UpdateAvatar(user, uc.AvatarDir, fullFileName)
}

func (uc *UserUseCase) GetUserByLogin(user string) (models.User, error) {
	return uc.Repository.GetUserByLogin(user)
}

func (uc *UserUseCase) GetProfileByLogin(login string) (models.User, error) {
	user, err := uc.Repository.GetUserByLogin(login)
	if err != nil {
		return models.User{}, err
	}
	return uc.GetOutputUserData(user), nil
}

func (uc *UserUseCase) Login(input models.UserSignIn) (models.User, error) {
	user, err := uc.GetUserByLogin(input.Login)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %v", err)
	}
	err = uc.CheckUserPassword(user.Password, input.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("wrong password: %v", err)
	}
	return user, nil
}

func (uc *UserUseCase) GetOutputUserData(user models.User) models.User {
	return models.User{
		Id:    user.Id,
		Name:  user.Name,
		Login: user.Login,
		Sex:   user.Sex,
		Image: user.Image,
		Email: user.Email,
		Theme: user.Theme,
	}
}

func (uc *UserUseCase) GetUserStat(id string) (models.UserStat, error) {
	return uc.Repository.GetUserStat(id)
}

func (uc *UserUseCase) CheckUserPassword(userPassword string, inputPassword string) error {
	return uc.Repository.CheckUserPassword(userPassword, inputPassword)
}
