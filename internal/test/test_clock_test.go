package test_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"test-project/internal/api"
	"test-project/internal/test"
)

func TestTestClock(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		now := time.Date(2025, 2, 5, 11, 42, 30, 0, time.UTC)

		test.SetMockClock(t, s, now)

		assert.Equal(t, now, s.Clock.Now())

		clock := test.GetMockClock(t, s.Clock)
		require.NotNil(t, clock)

		assert.Equal(t, now, clock.Now())
	})
}
