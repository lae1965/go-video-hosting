package service

import (
	"context"
	"database/sql"
	"fmt"
	grpcclient "go-video-hosting/gRPC/client"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"

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

func (userService *UserService) CreateUser(user model.Users) (*model.UserCreateResponse, *errors.AppError) {
	var appErr *errors.AppError
	transaction, err := userService.createTransaction()
	if err != nil {
		return nil, errors.New(errors.UnknownError, fmt.Sprintf("failed create transaction: %s", err.Error()))
	}

	defer func() {
		if err != nil || appErr != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	hash, err := userService.GenerateHashPassword(user.Password)
	if err != nil {
		return nil, errors.New(errors.UnknownError, fmt.Sprintf("failed generate hashPassword: %s", err.Error()))
	}
	user.Password = hash
	user.ActivateLink = uuid.New().String()

	userId, appErr := userService.dbUser.CreateUser(transaction, user)
	if appErr != nil {
		appErr.Message = fmt.Sprintf("failed saving new user: %s", appErr.Message)
		return nil, appErr
	}

	user.Id = userId

	tokenResponse, err := userService.token.CreateTokens(transaction, user, 0)
	if err != nil {
		return nil, errors.New(errors.UnknownError, fmt.Sprintf("failed creating tokens: %s", err.Error()))
	}

	err = SendMail(user.Email, user.ActivateLink)
	if err != nil {
		return nil, errors.New(errors.UnknownError, fmt.Sprintf("failed sending email message: %s", err.Error()))
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

func (userService *UserService) Login(user model.Users) (*model.UserResponse, *errors.AppError) {
	newUser, err := userService.dbUser.GetUserByEmail(user.Email)
	if err != nil {
		return nil, errors.New(errors.NotFound, fmt.Sprintf("user with email %s not found: %s", user.Email, err.Error()))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(user.Password)); err != nil {
		return nil, errors.New(errors.Unauthorization, fmt.Sprintf("wrong password: %s", err.Error()))
	}

	tokenResponse, err := userService.token.CreateTokens(nil, *newUser, 0)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err.Error())
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

func (userService *UserService) Refresh(refreshToken string) (*model.UserResponse, *errors.AppError) {
	userId, err := userService.token.ValidateToken(refreshToken, os.Getenv("REFRESH_KEY"))
	if err != nil {
		return nil, errors.New(errors.Unauthorization, err.Error())
	}

	refreshTokenId, err := userService.token.GetTokenIdByToken(refreshToken)
	if err != nil {
		return nil, errors.New(errors.Unauthorization, fmt.Sprintf("refreshtoken is not found in DB: %s", err.Error()))
	}

	userFromDB, err := userService.dbUser.GetUserForRefreshById(userId)
	if err != nil {
		return nil, errors.New(errors.Unauthorization, fmt.Sprintf("user with such refreshtoken is not exist: %s", err.Error()))
	}

	tokenResponse, err := userService.token.CreateTokens(nil, *userFromDB, refreshTokenId)
	if err != nil {
		return nil, errors.New(errors.Unauthorization, err.Error())
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

func (userService *UserService) SaveAvatar(id int, fileName string) *errors.AppError {
	oldFileName, appErr := userService.dbUser.GetAvatarByUserId(id)
	if appErr != nil {
		return errors.New(appErr.Type, fmt.Sprintf("can't get old avatarFileName: %s", appErr.Message))
	}

	newFileName, err := userService.grpcClient.SendToGRPCServer(context.Background(), fileName)
	if err != nil {
		return errors.New(errors.UnknownError, fmt.Sprintf("can't save file to gRPC-server: %s", err.Error()))
	}

	if err := userService.dbUser.UpdateUser(id, map[string]interface{}{"avatar": newFileName}); err != nil {
		userService.grpcClient.DeleteFromGRPCServer(context.Background(), newFileName)
		return errors.New(err.Type, fmt.Sprintf("can't save fileName to database: %s", err.Message))
	}

	if oldFileName != "" {
		userService.grpcClient.DeleteFromGRPCServer(context.Background(), oldFileName)
	}

	return nil
}

func (userService *UserService) GetAvatar(id int, sendChunk func(int64, string, []byte) error) *errors.AppError {
	avatarFileName, appErr := userService.dbUser.GetAvatarByUserId(id)
	if appErr != nil {
		appErr.Message = fmt.Sprintf("can't get avatarFileName: %s", appErr.Message)
		return appErr
	}

	if avatarFileName == "" {
		return errors.New(errors.EmptyField, "this user has no avatar")
	}

	if err := userService.grpcClient.GetFromGRPCServer(context.Background(), avatarFileName, sendChunk); err != nil {
		return errors.New(errors.UnknownError, fmt.Sprintf("can't get avatar: %s", err.Error()))
	}

	return nil
}

func (userService *UserService) DeleteAvatar(id int) *errors.AppError {
	avatarFileName, err := userService.dbUser.GetAvatarByUserId(id)
	if err != nil {
		return errors.New(err.Type, fmt.Sprintf("can't get avatarFileName: %s", err.Message))
	}

	if avatarFileName == "" {
		return errors.New(errors.EmptyField, "this user has no avatar")
	}

	if err := userService.dbUser.UpdateUser(id, map[string]interface{}{"avatar": ""}); err != nil {
		return errors.New(err.Type, fmt.Sprintf("can't delete avatarFileName from DB: %s", err.Message))
	}

	if err := userService.grpcClient.DeleteFromGRPCServer(context.Background(), avatarFileName); err != nil {
		return errors.New(errors.UnknownError, fmt.Sprintf("can't delete avatar from gRPC-server: %s", err.Error()))
	}

	return nil
}

func (userService *UserService) UpdateUser(id int, data map[string]interface{}) *errors.AppError {
	return userService.dbUser.UpdateUser(id, data)
}

func (userService *UserService) DeleteUser(id int) *errors.AppError {
	//TODO - удалить все видео user'а и аватар с gRPC - сервера
	return userService.dbUser.DeleteUser(id)
}

func (userService *UserService) Activate(activateLink string) *errors.AppError {
	userId, err := userService.dbUser.GetUserByActivateLink(activateLink)
	if err != nil {
		err.Message = fmt.Sprintf("can not find user by activate link: %s", err.Message)
		return err
	}

	if err := userService.dbUser.UpdateUser(userId, map[string]interface{}{"isActivate": true}); err != nil {
		err.Message = fmt.Sprintf("can not update field isActivate: %s", err.Message)
		return err
	}

	return nil
}

func (userService *UserService) GetAll() ([]*model.FindUsers, error) {
	return userService.dbUser.GetAll()
}

func (userService *UserService) GetById(id int) (*model.FindUsers, *errors.AppError) {
	return userService.dbUser.GetById(id)
}

func (userService *UserService) GetNickNameById(id int) (string, *errors.AppError) {
	return userService.dbUser.GetNickNameById(id)
}

func (userService *UserService) CheckIsNickNameEmailUnique(nickName string, email string) (bool, string, error) {
	notUniqueList := []string{}
	isUniqie := func(key string, value string) (bool, error) {
		var isValueUnique bool
		var err error

		if value == "" {
			isValueUnique = true
		} else {
			isValueUnique, err = userService.dbUser.CheckIsUnique(key, value)
			if err != nil {
				return false, err
			}
			if !isValueUnique {
				notUniqueList = append(notUniqueList, key)
			}
		}

		return isValueUnique, nil
	}

	isNickNameUnique, err := isUniqie("nickName", nickName)
	if err != nil {
		return false, "", err
	}

	isEmailUnique, err := isUniqie("email", email)
	if err != nil {
		return false, "", err
	}

	message := ""
	if len(notUniqueList) > 0 {
		message = fmt.Sprintf("%s not unique", strings.Join(notUniqueList, " and "))
	}

	return isEmailUnique && isNickNameUnique, message, nil
}

func (userService *UserService) ChangePassword(userId int, refreshTokenId int, oldPassword string, newPassword string) *errors.AppError {
	dbPassword, appErr := userService.dbUser.GetPasswordByUserId(userId)
	if appErr != nil {
		appErr.Message = fmt.Sprintf("wrong getting old Password from db: %s", appErr.Message)
		return appErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(oldPassword)); err != nil {
		return errors.New(errors.NotEqual, fmt.Sprintf("wrong old Password: %s", err.Error()))
	}

	hashPassword, err := userService.GenerateHashPassword(newPassword)
	if err != nil {
		return errors.New(errors.UnknownError, fmt.Sprintf("wrong generating hashPassword: %s", err.Error()))
	}

	if err := userService.dbUser.UpdateUser(userId, map[string]interface{}{"password": hashPassword}); err != nil {
		err.Message = fmt.Sprintf("wrong updating password in DB: %s", err.Message)
		return err
	}

	if err := userService.token.DeleteTokenFromOtherDevices(userId, refreshTokenId); err != nil {
		return errors.New(errors.UnknownError, fmt.Sprintf("wrong logouting user from other devices: %s", err.Error()))
	}

	return nil
}
