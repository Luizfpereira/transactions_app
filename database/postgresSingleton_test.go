package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

// TestConnectSingleton verifies if the gorm instance is nil at the beginning and if it is not after the first execution
// of the singleton. After multiple executions, the function should return the same address of the instance variable
func TestConnectSingleton(t *testing.T) {
	assert.Nil(t, instance)
	instanceSingleton := ConnectSingleton(sqlite.Open(":memory:"))
	assert.NotNil(t, instanceSingleton)
	assert.Equal(t, &instance, &instanceSingleton)
}
