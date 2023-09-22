package usecase

import (
	"log"
	"transactions_app/entity"
	"transactions_app/gateway"
)

type TransactionUsecase struct {
	gateway gateway.TransactionGateway
}

func NewTransactionUsecase(gateway gateway.TransactionGateway) *TransactionUsecase {
	return &TransactionUsecase{gateway: gateway}
}

func (t *TransactionUsecase) CreateTransaction(input entity.TransactionInput) (*entity.Transaction, error) {
	transaction, err := entity.NewTransaction(input.TransactionDate, input.Description, input.PurchaseAmount)
	if err != nil {
		return nil, err
	}

	transactionOutput, err := t.gateway.Create(transaction)
	if err != nil {
		return nil, err
	}
	return transactionOutput, nil
}

// GetTransactionByID
// GetTransactionsCurrency or currencies in plural

func (t *TransactionUsecase) GetTransactionsCurrency(currency string) ([]entity.TransactionConvertedOutput, error) {
	log.Println(t.gateway.GetTransactions())
	return nil, nil
}
