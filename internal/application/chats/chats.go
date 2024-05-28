package chats

import (
	"api/internal/core"
	"api/internal/infrastructure"
	"errors"
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

func (s *chats) AddMessage(message *core.Message) error {
	return s.repo.AddMessage(message)
}

func (s *chats) GetAll(userId int) (*[]core.Chat, error) {
	return s.repo.GetAll(userId)
}

func (s *chats) GetMessages(userId, chatId int) (*[]core.Message, error) {
	if userId == chatId {
		return nil, errors.New("userId == chatId")
	}
	if err := s.usersRepo.ExistsById(chatId); err != nil {
		return nil, err
	}

	return s.repo.GetMessages(chatId)
}
