package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"test-project/internal/util"
)

func TestFalseIfNil(t *testing.T) {
	b := true
	assert.True(t, util.FalseIfNil(&b))
	b = false
	assert.False(t, util.FalseIfNil(&b))
	assert.False(t, util.FalseIfNil(nil))
}
