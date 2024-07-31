package repository

import (
	"Messaggio/producer/internal/model"
)

func (r *Repository) SaveMessage(message *model.Message) (string, error) {
	query := `INSERT INTO messages (content, create_time) VALUES ($1, $2) RETURNING id`
	var id string
	err := r.conn.Get(&id, query, message.Content, message.CreateTime)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) GetMessage(id string) (*model.Message, error) {
	query := `SELECT id, content, is_processed, create_time FROM messages WHERE id = $1`

	message := make([]*model.Message, 0, 1)

	err := r.conn.Select(&message, query, id)

	if len(message) == 0 {
		return nil, model.ErrNotFound
	}

	if err != nil {
		return nil, err

	}

	return message[0], nil
}

func (r *Repository) GetMessages() ([]*model.Message, error) {
	query := `SELECT id, content, is_processed, create_time FROM messages`
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

func (r *Repository) DeleteMessage(id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	res, err := r.conn.Exec(query, id)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *Repository) GetStatsByHour() ([]*model.ByHours, error) {
	query := `SELECT DATE_TRUNC('hour', create_time) as hour, count(*) as count FROM messages GROUP BY hour`
	stats := make([]*model.ByHours, 0)
	err := r.conn.Select(&stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
