package repository

import (
	"Messaggio/consumer/internal/model"
	"fmt"
	"github.com/rs/zerolog/log"
)

func (r *Repository) EmptyProcess(message *model.Message) error {
	log.Info().Msg(message.ID + " " + message.Content)
	query := `UPDATE messages SET is_processed = true WHERE id = $1`

	res, err := r.conn.Exec(query, message.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		//TODO: make error in variable
		return fmt.Errorf("no rows affected")
	}

	return nil
}
