package repository

import (
	"context"
	"time"

	"Hackathon/internal/domain"
	"github.com/jackc/pgx/v5"
)

type ConversationRepo interface {
	GetRecords() ([]domain.Record, error)
	GetRecord(ID int) (domain.Record, error)
	InsertMainRecordInfo(audioName string, createdAt time.Time) (domain.Record, error)
	InsertAdditionRecordInfo(id int, text string, goodPercent int, badPercent int) error
}

type conversationRepo struct {
	conn *pgx.Conn
}

func NewConversationRepo(conn *pgx.Conn) ConversationRepo {
	return &conversationRepo{
		conn: conn,
	}
}

const getRecordsQuery = `SELECT conversation_id, audio_text, audio_name, created_at, good_percent,
bad_percent FROM conversation`

func (r *conversationRepo) GetRecords() ([]domain.Record, error) {
	rows, err := r.conn.Query(context.Background(), getRecordsQuery)
	if err != nil {
		return nil, err
	}

	records := make([]domain.Record, 0)
	for rows.Next() {
		record := domain.Record{}
		err = rows.Scan(
			&record.ID,
			&record.Text,
			&record.AudioName,
			&record.CreatedAt,
			&record.GoodPercent,
			&record.BadPercent,
		)

		if err != nil {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

const getRecordQuery = `SELECT conversation_id, audio_text, audio_name, created_at, good_percent,
bad_percent FROM conversation
WHERE conversation_id=$1`

func (r *conversationRepo) GetRecord(ID int) (domain.Record, error) {
	row := r.conn.QueryRow(context.Background(), getRecordQuery, ID)
	record := domain.Record{}

	err := row.Scan(
		&record.ID,
		&record.Text,
		&record.AudioName,
		&record.CreatedAt,
		&record.GoodPercent,
		&record.BadPercent,
	)

	return record, err
}

const insertMainRecordInfoQuery = `
INSERT INTO conversation(audio_name, created_at)
VALUES ($1, $2) RETURNING conversation_id`

func (r *conversationRepo) InsertMainRecordInfo(audioName string, createdAt time.Time) (domain.Record, error) {
	row := r.conn.QueryRow(context.Background(), insertMainRecordInfoQuery, audioName, createdAt)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return domain.Record{}, err
	}

	return domain.Record{
		ID:        id,
		AudioName: audioName,
		CreatedAt: createdAt,
	}, nil

}

const insertAdditionRecordInfoQuery = `
UPDATE conversation
SET audio_text=$1, good_percent=$2, bad_percent=$3
WHERE conversation_id=$4
`

func (r *conversationRepo) InsertAdditionRecordInfo(id int, text string, goodPercent int, badPercent int) error {
	_, err := r.conn.Exec(context.Background(), insertAdditionRecordInfoQuery, text, goodPercent, badPercent, id)
	return err
}
