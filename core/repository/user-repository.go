package repository

import (
	"database/sql"
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	errorutil "mfawzanid/warehouse-commerce/utils/error"

	"github.com/lib/pq"
)

type UserRepositoryInterface interface {
	GetUser(req *entity.GetUserRequest) (*entity.User, error)
	InsertUser(user *entity.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &userRepository{
		db: db,
	}
}

const (
	uniqueViolationErrorCode = "23505"
)

func (r *userRepository) GetUser(req *entity.GetUserRequest) (*entity.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	query := `SELECT id, email, phone_number FROM users`

	values := []interface{}{}
	if req.Email != "" {
		query += " WHERE email = $1"
		values = append(values, req.Email)
	} else if req.PhoneNumber != "" {
		query += " WHERE phone_number = $1"
		values = append(values, req.PhoneNumber)
	}

	user := &entity.User{}

	err := r.db.QueryRow(query, values...).Scan(&user.Id, &user.Email, &user.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorutil.NewErrorCode(errorutil.ErrNotFound, fmt.Errorf("error repo get user by identifier '%s' or '%s'", req.Email, req.PhoneNumber))
		} else {
			return nil, fmt.Errorf("error repo get user: %v", err.Error())
		}
	}

	return user, nil
}

func (r *userRepository) InsertUser(user *entity.User) error {
	query := "INSERT INTO users (id, email, phone_number) VALUES ($1, $2, $3)"

	_, err := r.db.Exec(query, user.Id, user.Email, user.PhoneNumber)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationErrorCode {
			return nil
		} else {
			return fmt.Errorf("error repo create user: %v", err.Error())
		}
	}

	return nil
}
