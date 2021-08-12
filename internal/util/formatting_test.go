package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAndSanitizeMsisdn(t *testing.T) {
	t.Run("invalid length", func(t *testing.T) {
		_, err := ValidateAndSanitizeMsisdn(" 62800000 ")
		assert.NotNil(t, err)
	})
	t.Run("exceed max length", func(t *testing.T) {
		_, err := ValidateAndSanitizeMsisdn(" 62800000000000000 ")
		assert.NotNil(t, err)
	})
	t.Run("trim space", func(t *testing.T) {
		msisdn, err := ValidateAndSanitizeMsisdn(" 628000000 ")
		assert.Nil(t, err)
		assert.Equal(t, "628000000", msisdn)
	})
	t.Run("start with zero", func(t *testing.T) {
		msisdn, err := ValidateAndSanitizeMsisdn("080000001")
		assert.Nil(t, err)
		assert.Equal(t, "6280000001", msisdn)
	})
}
