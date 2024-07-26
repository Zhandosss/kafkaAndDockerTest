package repository

import (
	"Messaggio/producer/internal/model"
)

func (r *Repository) SaveMessage(message *model.Message) (string, error) {
	query := `INSERT INTO messages (content) VALUES ($1) RETURNING id`
	var id string
	err := r.conn.Get(&id, query, message.Content)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) GetMessage(id string) (*model.Message, error) {
	query := `SELECT id, content, is_processed FROM messages WHERE id = $1`
	message := make([]*model.Message, 0, 1)
	err := r.conn.Select(&message, query, id)
	if err != nil {
		return nil, err

	}
	return message[0], nil
}

func (r *Repository) GetMessages() ([]*model.Message, error) {
	query := `SELECT id, content, is_processed FROM messages`
	messages := make([]*model.Message, 0)
	err := r.conn.Select(&messages, query)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *Repository) DeleteMessages() error {
	query := `DELETE FROM messages`
	_, err := r.conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
