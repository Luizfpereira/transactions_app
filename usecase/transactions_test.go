package usecase

import (
	"errors"
	"testing"
	"time"
	"transactions_app/entity"
	mock_gateway "transactions_app/gateway/mocks"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type TransactionsUsecaseSuite struct {
	suite.Suite
	mockRepo *mock_gateway.MockTransactionGateway
}

func TestTransactionsUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TransactionsUsecaseSuite))
}

func (t *TransactionsUsecaseSuite) SetupTest() {
	ctrl := gomock.NewController(t.T())
	mockTransactionRepo := mock_gateway.NewMockTransactionGateway(ctrl)
	t.mockRepo = mockTransactionRepo
}

func (t *TransactionsUsecaseSuite) TestCreateTransaction() {
	inputDate := "2012-08-08 10:00:00"
	layout := "2006-01-02 15:04:05"
	parsedTime, _ := time.Parse(layout, inputDate)
	input := entity.TransactionInput{
		Description:     "test",
		TransactionDate: parsedTime,
		PurchaseAmount:  decimal.New(0, 0),
	}

	transactionFromInput, err := entity.NewTransaction(input.TransactionDate, input.Description, input.PurchaseAmount)
	t.Nil(err)

	var res entity.Transaction
	res.ID = 1
	res.Description = transactionFromInput.Description
	res.PurchaseAmount = transactionFromInput.PurchaseAmount
	res.TransactionDate = input.TransactionDate

	t.mockRepo.EXPECT().Create(transactionFromInput).Return(&res, nil)

	u := NewTransactionUsecase(t.mockRepo)
	output, err := u.CreateTransaction(input)
	t.Nil(err)
	t.Equal(1, int(output.ID))
	t.Equal(output, &res)
}

func (t *TransactionsUsecaseSuite) TestTransactionByIdCurrencyEmpty() {
	outputDB := entity.Transaction{
		TransactionDate: time.Now(),
		Description:     "test",
		PurchaseAmount:  decimal.New(50, 0),
	}
	outputDB.ID = 1
	t.mockRepo.EXPECT().GetTransactionById(1).Return(&outputDB, nil)

	u := NewTransactionUsecase(t.mockRepo)
	converted, err := u.GetTransactionByIdCurrency("", 1)
	t.Nil(err)
	t.Equal(outputDB.PurchaseAmount, converted.PurchaseAmount)
	t.Nil(converted.ExchangeRate)
	t.Nil(converted.ConvertedAmount)
	t.Nil(converted.Error)
}

func (t *TransactionsUsecaseSuite) TestTransactionByIdCurrencyEmptyNoID() {
	errNotRegistered := errors.New("id not registered")
	t.mockRepo.EXPECT().GetTransactionById(1).Return(nil, errors.New("error retrieving id"))
	u := NewTransactionUsecase(t.mockRepo)
	converted, err := u.GetTransactionByIdCurrency("", 1)
	t.Nil(converted)
	t.EqualError(err, errNotRegistered.Error())
}

func (t *TransactionsUsecaseSuite) TestTransactionByIdCurrency() {
	date := time.Date(2020, 01, 01, 10, 0, 0, 0, time.UTC)
	exchangeRate := decimal.New(4475, -3)

	outputDB := entity.Transaction{
		TransactionDate: date,
		Description:     "test",
		PurchaseAmount:  decimal.New(100, 0),
	}
	outputDB.ID = 1

	convertedAmount := exchangeRate.Mul(outputDB.PurchaseAmount).Round(2)

	t.mockRepo.EXPECT().GetTransactionById(1).Return(&outputDB, nil)
	expected := entity.TransactionConvertedOutput{
		Id:              1,
		Description:     outputDB.Description,
		TransactionDate: outputDB.TransactionDate,
		PurchaseAmount:  outputDB.PurchaseAmount,
		ExchangeRate:    &exchangeRate,
		ConvertedAmount: &convertedAmount,
	}
	u := NewTransactionUsecase(t.mockRepo)
	converted, err := u.GetTransactionByIdCurrency("Brazil-Real", 1)
	t.Nil(err)
	t.Equal(expected.PurchaseAmount, converted.PurchaseAmount)
	t.Equal(*expected.ExchangeRate, *converted.ExchangeRate)
	t.Equal(*expected.ConvertedAmount, *converted.ConvertedAmount)
	t.Nil(converted.Error)
}

func (t *TransactionsUsecaseSuite) TestTransactionByIdWrongCurrency() {
	date := time.Date(2020, 01, 01, 10, 0, 0, 0, time.UTC)

	outputDB := entity.Transaction{
		TransactionDate: date,
		Description:     "test",
		PurchaseAmount:  decimal.New(100, 0),
	}

	t.mockRepo.EXPECT().GetTransactionById(1).Return(&outputDB, nil)

	u := NewTransactionUsecase(t.mockRepo)
	converted, err := u.GetTransactionByIdCurrency("WrongBrazil-Real", 1)
	t.NotNil(err)
	t.Nil(converted)
}

func (t *TransactionsUsecaseSuite) TestTransactionsCurrencyEmpty() {
	date := time.Date(2020, 01, 01, 10, 0, 0, 0, time.UTC)
	outputDB := []entity.Transaction{
		{
			TransactionDate: date,
			Description:     "test",
			PurchaseAmount:  decimal.New(1, 0),
		},
		{
			TransactionDate: date,
			Description:     "test2",
			PurchaseAmount:  decimal.New(50, 0),
		},
	}
	t.mockRepo.EXPECT().GetTransactions().Return(outputDB, nil)

	u := NewTransactionUsecase(t.mockRepo)
	converted, err := u.GetTransactionsCurrency("")
	t.Nil(err)

	for i := range converted {
		t.Equal(outputDB[i].PurchaseAmount, converted[i].PurchaseAmount)
		t.Equal(outputDB[i].Description, converted[i].Description)
		t.Nil(converted[i].ExchangeRate)
		t.Nil(converted[i].ConvertedAmount)
		t.Nil(converted[i].Error)
	}
}

func (t *TransactionsUsecaseSuite) TestTransactionsWrongCurrency() {
	date := time.Date(2020, 01, 01, 10, 0, 0, 0, time.UTC)

	outputDB := []entity.Transaction{
		{
			TransactionDate: date,
			Description:     "test",
			PurchaseAmount:  decimal.New(1, 0),
		},
		{
			TransactionDate: date,
			Description:     "test2",
			PurchaseAmount:  decimal.New(50, 0),
		},
	}

	t.mockRepo.EXPECT().GetTransactions().Return(outputDB, nil)

	u := NewTransactionUsecase(t.mockRepo)
	converted, err := u.GetTransactionsCurrency("WrongBrazil-Real")
	t.Nil(err)
	for _, c := range converted {
		t.Equal("purchase cannot be converted to target currency", *c.Error)
	}
}

func (t *TransactionsUsecaseSuite) TestTransactionsCurrency() {
	date := time.Date(2020, 01, 01, 10, 0, 0, 0, time.UTC)
	exchangeRate := decimal.New(4475, -3)

	outputDB := []entity.Transaction{
		{
			TransactionDate: date,
			Description:     "test",
			PurchaseAmount:  decimal.New(1, 0),
		},
		{
			TransactionDate: date,
			Description:     "test2",
			PurchaseAmount:  decimal.New(50, 0),
		},
	}
	t.mockRepo.EXPECT().GetTransactions().Return(outputDB, nil)
	var expected []entity.TransactionConvertedOutput

	for _, out := range outputDB {
		convertedAmount := exchangeRate.Mul(out.PurchaseAmount).Round(2)
		exp := entity.TransactionConvertedOutput{
			Description:     out.Description,
			TransactionDate: out.TransactionDate,
			PurchaseAmount:  out.PurchaseAmount,
			ExchangeRate:    &exchangeRate,
			ConvertedAmount: &convertedAmount,
		}
		expected = append(expected, exp)
	}
	u := NewTransactionUsecase(t.mockRepo)
	converted, err := u.GetTransactionsCurrency("Brazil-Real")
	t.Nil(err)
	for i, c := range converted {
		t.Equal(expected[i].PurchaseAmount, c.PurchaseAmount)
		t.Equal(*expected[i].ExchangeRate, *c.ExchangeRate)
		t.Equal(*expected[i].ConvertedAmount, *c.ConvertedAmount)
		t.Nil(c.Error)
	}
}
