package service

import (
	"fmt"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type CustomClaims struct {
	UserId   int    `json:"userId"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	dbToken database.Token
}

func NewTokenService(dbToken database.Token) *TokenService {
	return &TokenService{dbToken: dbToken}
}

func (s *TokenService) CreateTokens(transaction *sqlx.Tx, user model.Users, refreshTokenId int) (*model.TokenResponse, error) {
	createToken := func(claims CustomClaims, key string) (string, error) {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := tok.SignedString([]byte(key))
		if err != nil {
			return "", err
		}
		return token, nil
	}

	claims := &CustomClaims{
		UserId:   user.Id,
		Nickname: user.NickName,
		Email:    user.Email,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	accessToken, accessErr := createToken(*claims, os.Getenv("ACCESS_KEY"))
	if accessErr != nil {
		return nil, fmt.Errorf("error generating accessToken: %s", accessErr.Error())
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 60))
	refreshToken, refreshErr := createToken(*claims, os.Getenv("REFRESH_KEY"))
	if refreshErr != nil {
		return nil, fmt.Errorf("error generating refreshToken: %s", refreshErr.Error())
	}

	tokenId, err := s.saveRefreshTokenToDB(transaction, user.Id, refreshToken, refreshTokenId)
	if err != nil {
		return nil, fmt.Errorf("error saving refreshToken: %s", err.Error())
	}

	return &model.TokenResponse{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		RefreshTokenId: tokenId,
	}, nil
}

func (s *TokenService) saveRefreshTokenToDB(transaction *sqlx.Tx, userId int, refreshToken string, refreshTokenId int) (int, error) {
	var err error

	if refreshTokenId < 1 {
		refreshTokenId, err = s.dbToken.CreateToken(transaction, model.Token{
			Token:  refreshToken,
			UserId: userId,
		})
	} else {
		err = s.dbToken.UpdateToken(refreshTokenId, refreshToken)
	}
	if err != nil {
		return 0, err
	}

	return refreshTokenId, nil
}

func (s *TokenService) RemoveToken(tokenId int) error {
	err := s.dbToken.RemoveToken(tokenId)

	return err
}

func (s *TokenService) ValidateToken(tokenString string, tokenKey string) (int, error) {
	var claims CustomClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenKey), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token is not valid or expiried")
	}

	return claims.UserId, nil
}

func (s *TokenService) GetTokenIdByToken(token string) (int, error) {
	return s.dbToken.GetTokenIdByToken(token)
}

func (s *TokenService) DeleteTokenFromOtherDevices(userId int, refreshTokenId int) error {
	return s.dbToken.DeleteTokenFromOtherDevices(userId, refreshTokenId)
}
