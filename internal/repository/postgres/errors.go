package postgres

import "errors"

var (
	errStartingTx        = errors.New("error starting transaction")
	errCommitingTx       = errors.New("error commiting transaction")
	errCreatingTask      = errors.New("error creating task")
	errUpdatingLastRunAt = errors.New("error updating last run at")
)
