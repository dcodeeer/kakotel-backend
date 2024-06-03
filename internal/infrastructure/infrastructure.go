package infrastructure

import (
	"api/internal/core"
	"api/internal/infrastructure/chats"
	"api/internal/infrastructure/estates"
	"api/internal/infrastructure/users"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	Users   IUsers
	Estates IEstates
	Chats   IChats
}

type IChats interface {
	Add(user1, user2 int) (int, error)
	AddMessage(message *core.Message) (*core.Message, error)

	GetAll(userId int) (*[]core.Chat, error)
	GetMessages(chatId int) (*[]core.Message, error)

	IsChatMember(userId, chatId int) error
	GetChatIdByMembers(user1, user2 int) (int, error)
}

type IUsers interface {
	Add(user *core.User) (int, error)

	ExistsById(userId int) error

	GetOneById(userId int) (*core.User, error)
	GetByToken(token string) (*core.User, error)
	GetByEmail(email string) (*core.User, error)

	EmailExists(email string) error
	CreateToken(userId int) (string, error)

	Update(user *core.User) error
	UpdatePhoto(userId int, path string) error

	ChangePassword(userId int, password string) error

	SendRecoveryKey(email string) error
	GetUserIdByRecoveryKey(key string) (int, error)
	DeleteRecoveryKey(key string) error
	SetPasswordByUserId(userId int, password string) error
}

type IEstates interface {
	Add(estate *core.Estate) (int, error)

	AddAddress(address *core.Address) error

	AddTempImage(path string) (int, error)
	GetTempImages(ids []int) ([]string, error)

	GetCategories() ([]core.Category, error)
	GetAmenities() ([]core.Amenity, error)

	GetAll() (*[]core.Estate, error)
	GetOne(id int) (*core.Estate, error)
	// Remove(ctx context.Context, id int) error
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		Users:   users.New(db),
		Estates: estates.New(db),
		Chats:   chats.New(db),
	}
}
