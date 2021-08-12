package referral

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupReferHandler(t *testing.T) {
	assert.Panics(t, func() {
		SetupReferHandler(nil)
	})
	assert.NotPanics(t, func() {
		SetupReferHandler(&referUseCase{})
	})
}
