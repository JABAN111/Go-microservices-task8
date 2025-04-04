package util

import (
	"database/sql"
	"log/slog"
)

func CommitOrRollback(tx *sql.Tx, err *error, log *slog.Logger) {
	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error("Failed to rollback transaction", "rollback_error", rollbackErr)
		} else {
			log.Debug("Transaction rolled back due to error")
		}
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			log.Error("Failed to commit transaction", "commit_error", commitErr)
		} else {
			log.Debug("Transaction committed")
		}
	}
}
