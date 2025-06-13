package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"mfawzanid/warehouse-commerce/core/entity"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
)

type AuthUsecaseInterface interface {
	CreateToken(username string) (tokenString string, err error)
	VerifyToken(tokenString string) (userId string, err error)
}

type authUsecase struct {
}

func NewAuthUsecase() AuthUsecaseInterface {
	return &authUsecase{}
}

const (
	tokenSecret = "token_secret" // TODO improve by set as config
)

func (u *authUsecase) CreateToken(userId string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		entity.ContextUserId: userId,
		"exp":                time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err = token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error create token: %v", err.Error())
	}

	return tokenString, nil
}

func (u *authUsecase) VerifyToken(tokenString string) (userId string, err error) {
	token, err := jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", errorutil.NewErrorCode(errorutil.ErrUnauthorized, fmt.Errorf("error verify token: %v", err.Error()))
	}
	if !token.Valid {
		return "", errorutil.NewErrorCode(errorutil.ErrUnauthorized, errors.New("error verify token: token is invalid"))
	}

	userId, ok := token.Claims.(jwt.MapClaims)["userId"].(string)
	if !ok {
		return "", errorutil.NewErrorCode(errorutil.ErrUnauthorized, errors.New("error verify token: user id is not found in token"))
	}

	return userId, nil
}
