package repository

import (
	"transactions_app/entity"

	"gorm.io/gorm"
)

type TransactionRepositoryPsql struct {
	Instance *gorm.DB
}

func NewTransactionRepositoryPsql(instance *gorm.DB) *TransactionRepositoryPsql {
	return &TransactionRepositoryPsql{Instance: instance}
}

func (t *TransactionRepositoryPsql) Create(transaction *entity.Transaction) (*entity.Transaction, error) {
	return nil, nil
}
