package test

import (
	"testing"
	"time"

	"github.com/dropbox/godropbox/time2"
	"test-project/internal/api"
)

func GetMockClock(t *testing.T, clock time2.Clock) *time2.MockClock {
	t.Helper()
	mc, ok := clock.(*time2.MockClock)
	if !ok {
		t.Fatalf("invalid clock type, got %T, want *time2.MockClock", clock)
	}

	return mc
}

func SetMockClock(t *testing.T, s *api.Server, time time.Time) {
	mockClock := GetMockClock(t, s.Clock)

	mockClock.Set(time)
}
