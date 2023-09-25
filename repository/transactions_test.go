package repository

import (
	"testing"
	"time"
	"transactions_app/entity"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TransactionRepoSuite struct {
	db *gorm.DB
	suite.Suite
}

func TestTransactionRepoSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepoSuite))
}

func (t *TransactionRepoSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Suite.T().Fatal(err)
	}
	db.AutoMigrate(entity.Transaction{})
	t.db = db

	repo := NewTransactionRepositoryPsql(t.db)
	test1, _ := entity.NewTransaction(time.Now(), "test1", decimal.New(500, 0))
	resp1, err := repo.Create(test1)
	t.Nil(err)
	t.Equal(test1, resp1)

	test2, _ := entity.NewTransaction(time.Now(), "test2", decimal.New(1000, 0))
	resp2, err := repo.Create(test2)
	t.Nil(err)
	t.Equal(test2, resp2)
}

func (t *TransactionRepoSuite) TestCreate() {
	repo := NewTransactionRepositoryPsql(t.db)
	test, err := entity.NewTransaction(time.Now(), "test", decimal.New(1, 0))
	t.Nil(err)
	resp, err := repo.Create(test)
	t.Nil(err)
	t.Equal(3, int(resp.ID))
}

func (t *TransactionRepoSuite) TestGetTransactions() {
	repo := NewTransactionRepositoryPsql(t.db)
	resp, err := repo.GetTransactions()
	t.Nil(err)
	t.Equal(2, len(resp))
}

func (t *TransactionRepoSuite) TestGetTransactionById() {
	repo := NewTransactionRepositoryPsql(t.db)
	resp, err := repo.GetTransactionById(1)
	t.Nil(err)
	t.Equal(1, int(resp.ID))
	t.Equal("test1", resp.Description)

	resp3, err := repo.GetTransactionById(3)
	t.NotNil(err)
	t.Nil(resp3)
}
