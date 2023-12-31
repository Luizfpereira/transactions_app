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
	Description     string          `json:"description" binding:"required,max=50"`
	TransactionDate time.Time       `json:"transaction_date" binding:"required"`
	PurchaseAmount  decimal.Decimal `json:"purchase_amount" binding:"required"`
}

type TransactionConvertedOutput struct {
	Id              int              `json:"id,omitempty"`
	Description     string           `json:"description,omitempty"`
	TransactionDate time.Time        `json:"transaction_date,omitempty"`
	PurchaseAmount  decimal.Decimal  `json:"purchase_amount,omitempty"`
	ExchangeRate    *decimal.Decimal `json:"exchange_rate,omitempty"`
	ConvertedAmount *decimal.Decimal `json:"converted_amount,omitempty"`
	Error           *string          `json:"error,omitempty"`
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

	if len(t.Description) > 50 {
		return errors.New("description max length: 50")
	}
	return nil
}
