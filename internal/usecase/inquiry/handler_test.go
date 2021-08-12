package inquiry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupInquiringHandler(t *testing.T) {
	assert.Panics(t, func() {
		SetupInquiringHandler(nil)
	})
	assert.NotPanics(t, func() {
		SetupInquiringHandler(&inquiryUseCase{})
	})
}
