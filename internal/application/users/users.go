package users

import (
	"api/internal/core"
	"api/internal/infrastructure"
	"errors"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type users struct {
	repo infrastructure.IUsers
}

func New(repo infrastructure.IUsers) *users {
	return &users{repo: repo}
}

func (s *users) SignUp(email, password string) (string, error) {
	err := s.repo.EmailExists(email)
	if err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	user := &core.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	userId, err := s.repo.Add(user)
	if err != nil {
		return "", err
	}

	return s.repo.CreateToken(userId)
}

func (s *users) SignIn(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return s.repo.CreateToken(user.ID)
}

func (s *users) GetOneInfo(userId int) (map[string]interface{}, error) {
	output := make(map[string]interface{})

	user, err := s.repo.GetOneById(userId)
	if err != nil {
		return output, err
	}

	output["id"] = user.ID
	output["firstname"] = user.FirstName
	output["lastname"] = user.LastName
	output["photo"] = user.Photo
	output["last_seen"] = user.LastSeen

	return output, nil
}

func (s *users) GetOneById(userId int) (*core.User, error) {
	return s.repo.GetOneById(userId)
}

func (s *users) GetByToken(token string) (*core.User, error) {
	return s.repo.GetByToken(token)
}

func (s *users) Update(user *core.User) error {
	return s.repo.Update(user)
}

func (s *users) UpdateLastSeen(userId int) error {
	return s.repo.UpdateLastSeen(userId)
}

func (s *users) UpdatePhoto(userId int, bytes []byte) (string, error) {
	filename, err := core.CreateFile(bytes)
	if err != nil {
		log.Println("error 1")
		return "", err
	}

	user, err := s.repo.GetOneById(userId)
	if err != nil {
		log.Println("error 1")
		return "", err
	}

	err = s.repo.UpdatePhoto(userId, filename)
	if err != nil {
		log.Println("error 1")
		return "", err
	}

	var oldFile string

	if user.Photo != nil {
		oldFile = *user.Photo
		os.Remove("uploads/" + oldFile)
	}

	return filename, err
}

func (s *users) ChangePassword(userId int, oldPassword, newPassword string) error {
	user, err := s.repo.GetOneById(userId)
	if err != nil {
		return err
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err == nil {
		return s.repo.ChangePassword(userId, string(hashedNewPassword))
	} else {
		return errors.New("passwords not equals")
	}
}

func (s *users) SendRecoveryKey(email string) error {
	if err := s.repo.EmailExists(email); err != nil {
		return err
	}
	return s.repo.SendRecoveryKey(email)
}

func (s *users) ConfirmRecoveryKey(key, password string) (string, error) {
	userId, err := s.repo.GetUserIdByRecoveryKey(key)
	if err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	if err := s.repo.SetPasswordByUserId(userId, string(hashedPassword)); err != nil {
		return "", err
	}

	if err := s.repo.DeleteRecoveryKey(key); err != nil {
		return "", err
	}

	return s.repo.CreateToken(userId)
}
