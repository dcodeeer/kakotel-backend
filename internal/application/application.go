package application

import (
	"api/internal/application/chats"
	"api/internal/application/estates"
	"api/internal/application/users"
	"api/internal/core"
	"api/internal/infrastructure"
)

type UseCase struct {
	Users   IUsers
	Estates IEstates
	Chats   IChats
}

type IUsers interface {
	SignUp(email, password string) (string, error)
	SignIn(email, password string) (string, error)

	SendRecoveryKey(email string) error
	ConfirmRecoveryKey(key, password string) (string, error)

	GetOneInfo(userId int) (map[string]interface{}, error)
	GetOneById(userId int) (*core.User, error)
	GetByToken(token string) (*core.User, error)

	Update(user *core.User) error
	UpdatePhoto(userId int, bytes []byte) (image string, err error)
	ChangePassword(userId int, oldPassword, newPassword string) error
}

type IBooking interface {
	Add(estate *core.Estate) error
	UpdateStatus(estateId int) error
}

type IEstates interface {
	Add(estate *core.Estate) (int, error)

	GetTempImages(ids []int) ([]string, error)
	AddTempImage(bytes []byte) (int, error)
	GetAmenities() ([]core.Amenity, error)
	GetCategories() ([]core.Category, error)
	// Approve(estateId int) error
	GetAll() (*[]core.Estate, error)
	GetById(id int) (*core.Estate, error)
	// Remove(id int) error
}

type IChats interface {
	AddMessage(message *core.Message) error

	GetAll(userId int) (*[]core.Chat, error)
	GetMessages(userId, chatId int) (*[]core.Message, error)
}

func New(repo *infrastructure.Repo) *UseCase {
	return &UseCase{
		Users:   users.New(repo.Users),
		Estates: estates.New(repo.Estates),
		Chats:   chats.New(repo.Chats, repo.Users),
	}
}
