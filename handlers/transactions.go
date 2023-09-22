package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"transactions_app/config"
	"transactions_app/usecase"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type TransactionHandler struct {
	Usecase *usecase.TransactionUsecase
	config  config.Config
}

func NewTransactionHandler(usecase *usecase.TransactionUsecase, config config.Config) *TransactionHandler {
	return &TransactionHandler{Usecase: usecase, config: config}
}

func (t *TransactionHandler) Create(ctx *gin.Context) {
	description := ctx.PostForm("description")
	transactionDate := ctx.PostForm("transaction_date")
	purchaseAmount := ctx.PostForm("purchase_amount")

	if description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "description is empty"})
		return
	}

	if len(description) > 50 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "description max length: 50"})
		return
	}
	layout := "2006-01-02 15:04:05"

	date, err := time.Parse(layout, transactionDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": fmt.Sprintf("failed parsing date: %s", err.Error())})
		return
	}
	log.Println(date)

	value, err := decimal.NewFromString(purchaseAmount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	truncatedValue := value.Truncate(2)

	if truncatedValue.LessThanOrEqual(decimal.New(0, 0)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "purchase value must be greater than 0"})
		return
	}

}
