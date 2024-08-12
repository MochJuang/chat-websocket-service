package repository

type DatabaseTransactionRepository interface {
	CommitTransaction() error
	RollbackTransaction() error
}
