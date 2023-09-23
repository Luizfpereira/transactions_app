package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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
	transactions, err := t.gateway.GetTransactions()
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions registered")
	}

	var transactionsOutput []entity.TransactionConvertedOutput
	for _, transaction := range transactions {
		output := entity.TransactionConvertedOutput{
			Id:              int(transaction.ID),
			Description:     transaction.Description,
			TransactionDate: transaction.TransactionDate,
			PurchaseAmount:  transaction.PurchaseAmount,
		}
		if currency != "" {
			treasuryFiscalData, err := getExchangeRates(currency, output.TransactionDate)
			if err != nil || len(treasuryFiscalData) == 0 {
				errStr := "purchase cannot be converted to target currency"
				output.Error = &errStr
			} else {
				convertedAmount := output.PurchaseAmount.Mul(treasuryFiscalData[0].ExchageRate)
				output.ConvertedAmount = &convertedAmount
				output.ExchangeRate = &treasuryFiscalData[0].ExchageRate
			}

		}
		transactionsOutput = append(transactionsOutput, output)
	}

	return transactionsOutput, nil
}

func getExchangeRates(currency string, transactionDate time.Time) ([]entity.TreasuryFiscalData, error) {
	infLimit := transactionDate.AddDate(0, -6, 0).Format("2006-01-02")
	supLimit := transactionDate.Format("2006-01-02")

	// The apiUrl returns only the fields country_currency_desc, exchange_rate, record_date and effective_date.
	// It seems that for some countries, different exchange rates are registered on the same record date, in which
	// one of them will be only effective in a future date. So, it is important to sort and return at first place the
	// most updated record and effective dates and then verify if they are in the specified range. The request also
	// filters for a specific currency
	apiUrl := "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange" +
		"?fields=country_currency_desc,exchange_rate,record_date,effective_date&sort=-record_date,-effective_date" +
		"&filter=record_date:gte:%s,record_date:lte:%s,effective_date:lte:%s,country_currency_desc:in:(%s)"
	url := fmt.Sprintf(apiUrl, infLimit, supLimit, supLimit, currency)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result entity.TreasuryApi
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}
