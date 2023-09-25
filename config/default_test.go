package config

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLoadConfig tests two cases: Loading a non existant file and an existant one
func TestLoadConfig(t *testing.T) {
	t.Run("testing non existant file", func(t *testing.T) {
		_, err := LoadConfig(".")
		log.Println(err)
		assert.NotNil(t, err)
	})
	t.Run("testing app.env config file", func(t *testing.T) {
		_, err := LoadConfig("../.")
		assert.Nil(t, err, "the error should be nil")
	})

}
