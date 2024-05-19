package service

import (
	"time"

	"Hackathon/internal/domain"
	"Hackathon/internal/repository"
)

type ConversationService interface {
	GetRecords() ([]domain.Record, error)
	GetRecord(ID int) (domain.Record, error)
	InsertMainRecordInfo(audioName string) (domain.Record, error)
	InsertAdditionRecordInfo(id int, text string, goodPercent int, badPercent int) error
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

func (s *conversationService) InsertMainRecordInfo(audioName string) (domain.Record, error) {
	createdAt := time.Now()

	return s.repo.InsertMainRecordInfo(audioName, createdAt)
}

func (s *conversationService) InsertAdditionRecordInfo(id int, text string, goodPercent int, badPercent int) error {
	return s.repo.InsertAdditionRecordInfo(id, text, goodPercent, badPercent)
}
