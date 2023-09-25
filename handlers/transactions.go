package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"transactions_app/config"
	"transactions_app/entity"
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

	value, err := decimal.NewFromString(purchaseAmount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	truncatedValue := value.Round(2)

	// if truncatedValue.LessThanOrEqual(decimal.New(0, 0)) {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "purchase value must be greater than 0"})
	// 	return
	// }

	input := entity.TransactionInput{
		Description:     description,
		TransactionDate: date.UTC(),
		PurchaseAmount:  truncatedValue,
	}

	output, err := t.Usecase.CreateTransaction(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": output})
}

func (t *TransactionHandler) GetTransactionsCurrency(ctx *gin.Context) {
	currency := ctx.Query("currency")
	output, err := t.Usecase.GetTransactionsCurrency(currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": output})
}

func (t *TransactionHandler) GetTransactionByIdCurrency(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "wrong id value - should be an integer"})
		return
	}
	currency := ctx.Query("currency")
	output, err := t.Usecase.GetTransactionByIdCurrency(currency, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": output})
}
