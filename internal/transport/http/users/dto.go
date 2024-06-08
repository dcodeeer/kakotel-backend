package users

import (
	"api/internal/core"
	"errors"
)

type signUpDTO struct {
	Fullname string `json:"fullname"`
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
	Fullname    string `json:"fullname"`
	Description string `json:"description"`
	DateOfBirth string `json:"date_of_birth"`
}

func (dto *updateDTO) ToUser() *core.User {
	return &core.User{
		Fullname:    &dto.Fullname,
		Description: &dto.Description,
		DateOfBirth: &dto.DateOfBirth,
	}
}

type changePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
