package entity

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	TransactionDate time.Time       `gorm:"not null"`
	Description     string          `gorm:"size:50; not null"`
	PurchaseAmount  decimal.Decimal `gorm:"type:numeric;not null"`
}

type TransactionInput struct {
	Description     string    `json:"description" binding:"required,max=50"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
	PurchaseAmount  string    `json:"purchase_amount" binding:"required"`
}

// NewTransaction validates and creates a new transaction instance
func NewTransaction(transactionDate time.Time, description string, purchaseAmount decimal.Decimal) (*Transaction, error) {
	transaction := &Transaction{
		TransactionDate: transactionDate,
		Description:     description,
		PurchaseAmount:  purchaseAmount,
	}
	if err := transaction.Validate(); err != nil {
		return nil, err
	}
	return transaction, nil
}

// Validate verifies if transaction date, description and purchase amount are valid
func (t *Transaction) Validate() error {
	if t.TransactionDate.IsZero() {
		return errors.New("transaction date is zero")
	}

	if t.Description == "" {
		return errors.New("description is empty")
	}

	// if t.PurchaseAmount == 0 {
	// 	return errors.New("purchase amount is zero")
	// }
	return nil
}
