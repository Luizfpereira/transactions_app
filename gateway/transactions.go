package gateway

import "transactions_app/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) (*entity.Transaction, error)
}
