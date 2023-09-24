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
	res := t.Instance.Create(&transaction)
	if res.Error != nil {
		return nil, res.Error
	}
	return transaction, nil
}

func (t *TransactionRepositoryPsql) GetTransactions() ([]entity.Transaction, error) {
	var transactionSlice []entity.Transaction
	res := t.Instance.Find(&transactionSlice)
	if res.Error != nil {
		return nil, res.Error
	}
	return transactionSlice, nil
}

func (t *TransactionRepositoryPsql) GetTransactionById(id int) (*entity.Transaction, error) {
	var transaction *entity.Transaction
	res := t.Instance.First(&transaction, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return transaction, nil
}
