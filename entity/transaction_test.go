package entity

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	t.Run("testing transaction with zero date", func(t *testing.T) {
		test1, err := NewTransaction(time.Time{}, "test", decimal.New(10, 0))
		assert.Nil(t, test1)
		assert.NotNil(t, err)
	})

	t.Run("testing transaction with empty descriptions", func(t *testing.T) {
		test1, err := NewTransaction(time.Now(), "", decimal.New(10, 0))
		assert.Nil(t, test1)
		assert.NotNil(t, err)
	})

	t.Run("testing transaction with length > 50", func(t *testing.T) {
		str := "a"
		for i := 0; i < 51; i++ {
			str += "a"
		}
		test1, err := NewTransaction(time.Now(), str, decimal.New(10, 0))
		assert.Nil(t, test1)
		assert.NotNil(t, err)
	})

	t.Run("testing transaction", func(t *testing.T) {
		test1, err := NewTransaction(time.Now(), "test", decimal.New(10, 0))
		assert.Nil(t, err)
		assert.NotNil(t, test1)
		assert.Equal(t, "test", test1.Description)
	})
}
