package service

import (
	"database/sql"
	"fmt"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	dbToken database.Token
}

func NewTokenService(dbToken database.Token) *TokenService {
	return &TokenService{dbToken: dbToken}
}

func (tokenService *TokenService) CreateTokens(transaction *sql.Tx, user model.Users) (*model.TokenResponse, error) {
	createToken := func(climes jwt.Claims, key string) (string, error) {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, climes)
		token, err := tok.SignedString([]byte(key))
		if err != nil {
			return "", err
		}
		return token, nil
	}

	climes := jwt.MapClaims{
		"id":       user.Id,
		"nickName": user.NickName,
		"email":    user.Email,
		"role":     "user",
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	accessToken, accessErr := createToken(climes, os.Getenv("ACCESS_KEY"))
	if accessErr != nil {
		return nil, fmt.Errorf("error generating accessToken: %s", accessErr.Error())
	}

	climes["exp"] = time.Now().Add(time.Hour * 24 * 60).Unix()
	refreshToken, refreshErr := createToken(climes, os.Getenv("REFRESH_KEY"))
	if refreshErr != nil {
		return nil, fmt.Errorf("error generating refreshToken: %s", refreshErr.Error())
	}

	tokenId, err := tokenService.dbToken.CreateToken(
		transaction,
		model.Token{
			Token:  refreshToken,
			UserId: user.Id,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error saving refreshToken: %s", err.Error())
	}

	return &model.TokenResponse{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		RefreshTokenId: tokenId,
	}, nil
}
