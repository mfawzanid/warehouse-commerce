package usecase

import (
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	"mfawzanid/warehouse-commerce/core/repository"
	serialutil "mfawzanid/warehouse-commerce/utils/serial"
)

type UserUsecaseInterface interface {
	RegisterUser(req *entity.RegisterUserRequest) (token string, err error)
	Login(req *entity.LoginRequest) (token string, err error)
}

type userUsecase struct {
	userRepo    repository.UserRepositoryInterface
	authUsecase AuthUsecaseInterface
}

const (
	userPrefixSerial = "USR"
)

func NewUserUsecase(userRepo repository.UserRepositoryInterface, authUsecase AuthUsecaseInterface) UserUsecaseInterface {
	return &userUsecase{userRepo, authUsecase}
}

func (u *userUsecase) RegisterUser(req *entity.RegisterUserRequest) (token string, err error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	userId, err := serialutil.GenerateId(userPrefixSerial)
	if err != nil {
		return "", fmt.Errorf("error register user in generating uuid: %v", err.Error())
	}

	var email, phoneNumber string
	if req.IdentifierType == entity.IdentifierTypeEmail {
		email = req.Identifier
	} else if req.IdentifierType == entity.IdentifierTypePhoneNumber {
		phoneNumber = req.Identifier
	}

	err = u.userRepo.InsertUser(&entity.User{
		Id:          userId,
		Email:       email,
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		return "", err
	}

	token, err = u.authUsecase.CreateToken(userId)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userUsecase) Login(req *entity.LoginRequest) (token string, err error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	getUserReq := &entity.GetUserRequest{}
	if req.IdentifierType == entity.IdentifierTypeEmail {
		getUserReq.Email = req.Identifier
	} else if req.IdentifierType == entity.IdentifierTypePhoneNumber {
		getUserReq.PhoneNumber = req.Identifier
	}

	user, err := u.userRepo.GetUser(getUserReq)
	if err != nil {
		return "", err
	}

	token, err = u.authUsecase.CreateToken(user.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}
