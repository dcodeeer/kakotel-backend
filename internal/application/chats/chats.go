package chats

import (
	"api/internal/core"
	"api/internal/infrastructure"
)

type chats struct {
	repo      infrastructure.IChats
	usersRepo infrastructure.IUsers
}

func New(repo infrastructure.IChats, usersRepo infrastructure.IUsers) *chats {
	return &chats{
		repo:      repo,
		usersRepo: usersRepo,
	}
}

func (s *chats) Add(user1, user2 int) (int, error) {
	return s.repo.Add(user1, user2)
}

func (s *chats) GetChatIdByMembers(user1, user2 int) (int, error) {
	return s.repo.GetChatIdByMembers(user1, user2)
}
func (s *chats) AddMessage(message *core.Message) (*core.Message, error) {
	return s.repo.AddMessage(message)
}

func (s *chats) GetAll(userId int) (*[]core.Chat, error) {
	return s.repo.GetAll(userId)
}

func (s *chats) GetMessages(userId, chatId int) (*[]core.Message, error) {
	return s.repo.GetMessages(chatId)
}
