package database

import (
	"fmt"
	"log"
	"testing"
	"transactions_app/entity"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestMigrate verifies if AutoMigrate returns an error and if it really creates a transactions table when executed
func TestMigrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("error opening gorm instance: %s", err.Error()))
	}
	err = db.AutoMigrate(entity.Transaction{})
	assert.Nil(t, err)

	rows, err := db.Table("sqlite_master").Select("name").Where("type = ? and name=?", "table", "transactions").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var tables []string
	var name string
	for rows.Next() {
		rows.Scan(&name)
		log.Println(name)
		tables = append(tables, name)
	}

	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "transactions", tables[0])
}
