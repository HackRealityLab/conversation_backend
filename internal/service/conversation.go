package service

import (
	"Hackathon/internal/domain"
	"Hackathon/internal/repository"
)

type ConversationService interface {
	GetRecords() ([]domain.Record, error)
	GetRecord(ID int) (domain.Record, error)
}

type conversationService struct {
	repo repository.ConversationRepo
}

func NewConversationService(repo repository.ConversationRepo) ConversationService {
	return &conversationService{
		repo: repo,
	}
}

func (s *conversationService) GetRecords() ([]domain.Record, error) {
	return s.repo.GetRecords()
}

func (s *conversationService) GetRecord(ID int) (domain.Record, error) {
	return s.repo.GetRecord(ID)
}
