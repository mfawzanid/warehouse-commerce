package entity

import (
	"fmt"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
)

const (
	ContextUserId = "userId"
)

// TODO: improve by using enum
const (
	IdentifierTypeEmail       = "email"
	IdentifierTypePhoneNumber = "phoneNumber"
)

type RegisterUserRequest struct {
	IdentifierType string // "email" or "phoneNumber"
	Identifier     string // email or phoneNumber value
}

func (r *RegisterUserRequest) Validate() error {
	if r.IdentifierType != IdentifierTypeEmail && r.IdentifierType != IdentifierTypePhoneNumber {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error register user validation: identifier type should be 'email' or 'phoneNumber'"))
	}
	if r.Identifier == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error register user validation: identifier is mandatory"))
	}
	return nil
}

type User struct {
	Id          string
	Email       string
	PhoneNumber string
}

type LoginRequest struct {
	IdentifierType string // "email" or "phoneNumber"
	Identifier     string // email or phoneNumber value
}

func (r *LoginRequest) Validate() error {
	if r.IdentifierType != IdentifierTypeEmail && r.IdentifierType != IdentifierTypePhoneNumber {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error login: identifier type should be 'email' or 'phoneNumber'"))
	}
	if r.Identifier == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error login: identifier is mandatory"))
	}
	return nil
}

type GetUserRequest struct {
	Id          string
	Email       string
	PhoneNumber string
}

func (r *GetUserRequest) Validate() error {
	if r.Email == "" && r.PhoneNumber == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error get user request validation: email or phone number is mandatory"))
	}
	return nil
}
