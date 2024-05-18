package repository

import (
	"context"

	"Hackathon/internal/domain"
	"github.com/jackc/pgx/v5"
)

type ConversationRepo interface {
	GetRecords() ([]domain.Record, error)
	GetRecord(ID int) (domain.Record, error)
}

type conversationRepo struct {
	conn *pgx.Conn
}

func NewConversationRepo(conn *pgx.Conn) {

}

const getRecordsQuery = `SELECT conversation_id, text, audio_name, created_at, good_percent,
is_ok FROM conversation`

func (r *conversationRepo) GetRecords() ([]domain.Record, error) {
	rows, err := r.conn.Query(context.Background(), getRecordsQuery)
	if err != nil {
		return nil, err
	}

	records := make([]domain.Record, 0)
	for rows.Next() {
		var record domain.Record
		err = rows.Scan(
			&record.ID,
			&record.Text,
			&record.AudioName,
			&record.CreatedAt,
			&record.GoodPercent,
			&record.IsOk,
		)

		if err != nil {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

const getRecordQuery = `SELECT conversation_id, text, audio_name, created_at, good_percent,
is_ok FROM conversation
WHERE conversation_id=$1`

func (r *conversationRepo) GetRecord(ID int) (domain.Record, error) {
	row := r.conn.QueryRow(context.Background(), getRecordQuery, ID)
	var record domain.Record

	err := row.Scan(
		&record.ID,
		&record.Text,
		&record.AudioName,
		&record.CreatedAt,
		&record.GoodPercent,
		&record.IsOk,
	)

	return record, err
}