package repository

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func (r *Repository) EmptyProcess(id string) error {
	log.Info().Msgf("message id is %s", id)
	query := `UPDATE messages SET is_processed = true WHERE id = $1`

	res, err := r.conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
