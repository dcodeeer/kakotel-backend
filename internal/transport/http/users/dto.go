package users

import (
	"api/internal/core"
	"errors"
)

type signUpDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *signUpDTO) Validate() error {
	if dto.Email == "" || dto.Password == "" {
		return errors.New("empty")
	}
	return nil
}

type updateDTO struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	DateOfBirth string `json:"date_of_birth"`
}

func (dto *updateDTO) ToUser() *core.User {
	return &core.User{
		FirstName:   &dto.FirstName,
		LastName:    &dto.LastName,
		DateOfBirth: &dto.DateOfBirth,
	}
}

type changePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
