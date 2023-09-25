package usecase

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ----- create Transaction
// deve receber um input e retornar toda a estrutura da transaction

// --------- Get Transactions currency
// caso 1: sem currency
// deve retornar todas as transactions armazenadas no mock - verificar id e length e validar se nao há campos no json de conversao
// caso 2: com currency
// deve retornar todas as transactions armazenadas no mock com os valores convertidos - validar length e valores

// --------- Get Transactions By Id currency
// caso 1: sem currency
// deve retornar uma transaction armazenada no mock - verificar id e validar se nao há campos no json de conversao
// caso 2: com currency
// deve retornar a transaction armazenadas no mock com os valores convertidos - validar length e valores

type TransactionsUsecaseSuite struct {
	suite.Suite
}

func TestTransactionsUsecaseSuite(t *testing.T) {

}
