package service

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"

	"os"
	"strings"

	"github.com/go-gomail/gomail"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type CallbackFunc func() (*sql.Tx, error)

type UserService struct {
	dbUser            database.Users
	token             Token
	createTransaction CallbackFunc
}

func NewUserService(dbUser database.Users, token Token, CreateTransaction CallbackFunc) *UserService {
	return &UserService{
		dbUser:            dbUser,
		token:             token,
		createTransaction: CreateTransaction,
	}
}

func (userService *UserService) CreateUser(user model.Users) (*model.UserResponse, error) {
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

	user.Password = userService.GenerateHashPassword(user.Password)
	user.ActivateLink = uuid.New().String()

	userId, err := userService.dbUser.CreateUser(transaction, user)
	if err != nil {
		return nil, err
	}

	user.Id = userId

	tokenResponse, err := userService.token.CreateTokens(transaction, user)
	if err != nil {
		return nil, err
	}

	err = SendMail(user.Email, user.ActivateLink)
	if err != nil {
		return nil, fmt.Errorf("error sending email message: %s", err.Error())
	}

	return &model.UserResponse{
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

func (userService *UserService) GenerateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}
