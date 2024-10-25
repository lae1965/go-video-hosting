package service

import (
	"context"
	"database/sql"
	"fmt"
	grpcclient "go-video-hosting/gRPC/client"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"

	"os"
	"strings"

	"github.com/go-gomail/gomail"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type CallbackFunc func() (*sql.Tx, error)

type UserService struct {
	dbUser            database.Users
	token             Token
	grpcClient        grpcclient.FilesGRPCClient
	createTransaction CallbackFunc
}

func NewUserService(dbUser database.Users, token Token, CreateTransaction CallbackFunc, grpcClient grpcclient.FilesGRPCClient) *UserService {
	return &UserService{
		dbUser:            dbUser,
		token:             token,
		grpcClient:        grpcClient,
		createTransaction: CreateTransaction,
	}
}

func (userService *UserService) CreateUser(user model.Users) (*model.UserCreateResponse, error) {
	transaction, err := userService.createTransaction()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	hash, err := userService.GenerateHashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %s", err.Error())
	}
	user.Password = hash
	user.ActivateLink = uuid.New().String()

	userId, err := userService.dbUser.CreateUser(transaction, user)
	if err != nil {
		return nil, err
	}

	user.Id = userId

	tokenResponse, err := userService.token.CreateTokens(transaction, user, 0)
	if err != nil {
		return nil, err
	}

	err = SendMail(user.Email, user.ActivateLink)
	if err != nil {
		return nil, fmt.Errorf("error sending email message: %s", err.Error())
	}

	return &model.UserCreateResponse{
		TokenResponse: tokenResponse,
		UserId:        userId,
	}, nil
}

func SendMail(email string, activateLink string) error {
	emailArray := strings.Split(email, "@")
	parsedEmail := strings.Split(emailArray[0], "+")[0]
	realEmail := fmt.Sprintf("%s@%s", parsedEmail, emailArray[1])

	link := fmt.Sprintf("%s:%s/api/user/activate/%s", viper.GetString("host"), viper.GetString("port"), activateLink)

	msg := gomail.NewMessage()
	msg.SetHeader("From", viper.GetString("mail.from"))
	msg.SetHeader("To", realEmail)
	msg.SetHeader("Subject", "Активация аккаунта на yo-tube.ru")
	msg.SetBody("text/html", fmt.Sprintf("<h2>Для активации аккаунта перейдите <a href=\"%s\">по ссылке</a>.</h2>", link))

	dialer := gomail.NewDialer(viper.GetString("mail.host"), viper.GetInt("mail.port"), "yo-tube", os.Getenv("MAIL_PASSWORD"))
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func (userService *UserService) GenerateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (userService *UserService) Login(user model.Users) (*model.UserResponse, error) {
	newUser, err := userService.dbUser.GetUserByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("user with email %s not found: %s", user.Email, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(user.Password))
	if err != nil {
		return nil, fmt.Errorf("wrong password: %s", err.Error())
	}

	tokenResponse, err := userService.token.CreateTokens(nil, *newUser, 0)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		TokenResponse: tokenResponse,
		Users:         newUser,
	}, nil
}

func (userService *UserService) Logout(refreshTokenId int) error {
	err := userService.token.RemoveToken(refreshTokenId)

	if err != nil {
		return fmt.Errorf("error removing token: %s", err.Error())
	}

	return nil
}

func (userService *UserService) Refresh(refreshToken string) (*model.UserResponse, error) {
	userId, err := userService.token.ValidateToken(refreshToken, os.Getenv("REFRESH_KEY"))
	if err != nil {
		return nil, err
	}

	refreshTokenId, err := userService.token.GetTokenIdByToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refreshtoken is not found in DB: %s", err.Error())
	}

	userFromDB, err := userService.dbUser.GetUserById(userId)
	if err != nil {
		return nil, fmt.Errorf("user with such refreshtoken is not exist: %s", err.Error())
	}

	tokenResponse, err := userService.token.CreateTokens(nil, *userFromDB, refreshTokenId)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		TokenResponse: tokenResponse,
		Users: &model.Users{
			Id:       userId,
			NickName: userFromDB.NickName,
			IsBanned: userFromDB.IsBanned,
			Role:     userFromDB.Role,
		},
	}, nil
}

func (userService *UserService) SaveAvatar(id int, fileName string) error {
	oldFileName, err := userService.dbUser.GetAvatarById(id)

	if err != nil {
		return fmt.Errorf("can't get avatarFileName: %s", err.Error())
	}

	newFileName, err := userService.grpcClient.SendToGRPCServer(context.Background(), fileName)
	if err != nil {
		return fmt.Errorf("can't save file to gRPC-server: %s", err.Error())
	}

	if err := userService.dbUser.UpdateAvatar(id, newFileName); err != nil {
		userService.grpcClient.DeleteFromGRPCServer(context.Background(), newFileName)
		return fmt.Errorf("can't save fileName to database: %s", err.Error())
	}

	if oldFileName != "" {
		userService.grpcClient.DeleteFromGRPCServer(context.Background(), oldFileName)
	}

	return nil
}

func (userService *UserService) GetAvatar(id int, sendChunk func(int64, string, []byte) error) error {
	avatarFileName, err := userService.dbUser.GetAvatarById(id)
	if err != nil {
		return fmt.Errorf("can't get avatarFileName: %s", err.Error())
	}

	if avatarFileName == "" {
		return fmt.Errorf("this user has no avatar")
	}

	if err := userService.grpcClient.GetFromGRPCServer(context.Background(), avatarFileName, sendChunk); err != nil {
		return fmt.Errorf("can't get avatar: %s", err.Error())
	}

	return nil
}

func (userService *UserService) DeleteAvatar(id int) error {
	avatarFileName, err := userService.dbUser.GetAvatarById(id)
	if err != nil {
		return fmt.Errorf("can't get avatarFileName: %s", err.Error())
	}

	if avatarFileName == "" {
		return fmt.Errorf("this user has no avatar")
	}

	if err := userService.grpcClient.DeleteFromGRPCServer(context.Background(), avatarFileName); err != nil {
		return fmt.Errorf("can't delete avatar from gRPC-server: %s", err.Error())
	}

	userService.dbUser.UpdateAvatar(id, "")

	return nil
}
