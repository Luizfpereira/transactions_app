package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
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

func (t *TransactionUsecase) GetTransactionsCurrency(currency string) ([]entity.TransactionConvertedOutput, error) {
	transactions, err := t.gateway.GetTransactions()
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions registered")
	}

	type Output struct {
		TransactionConverted *entity.TransactionConvertedOutput
		Error                error
	}

	// poolSize refers to the number of goroutines that could be created at the same time. This number is limited so the app
	// would not run out of memory resources in case they were created accordingly to a large number of transactions
	// existing in the database
	poolSize := 10
	queue := make(chan entity.Transaction)
	output := make(chan Output)

	var wg sync.WaitGroup
	wg.Add(len(transactions))

	go func(queue chan<- entity.Transaction) {
		for _, transaction := range transactions {
			queue <- transaction
		}
		close(queue)
	}(queue)

	for i := 0; i < poolSize; i++ {
		go func(queue <-chan entity.Transaction, output chan<- Output) {
			for t := range queue {
				defer wg.Done()
				var out Output
				converted := entity.TransactionConvertedOutput{
					Id:              int(t.ID),
					Description:     t.Description,
					TransactionDate: t.TransactionDate,
					PurchaseAmount:  t.PurchaseAmount,
				}
				out.TransactionConverted = &converted
				if currency != "" {
					treasuryFiscalData, err := getExchangeRates(currency, t.TransactionDate)
					if err != nil {
						out.Error = err
						output <- out
						return
					} else if len(treasuryFiscalData) == 0 {
						errStr := "purchase cannot be converted to target currency"
						converted.Error = &errStr
					} else {
						convertedAmount := t.PurchaseAmount.Mul(treasuryFiscalData[0].ExchageRate)
						converted.ConvertedAmount = &convertedAmount
						converted.ExchangeRate = &treasuryFiscalData[0].ExchageRate
					}
				}
				output <- out
			}
		}(queue, output)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	var transactionsOutput []entity.TransactionConvertedOutput
	for t := range output {
		if t.Error != nil {
			err = t.Error
		}
		transactionsOutput = append(transactionsOutput, *t.TransactionConverted)
	}
	if err != nil {
		return nil, err
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
		return nil, errors.New("currency could not be retrieved")
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
