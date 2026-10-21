package golf

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
)

type DynamicValue func(t *testing.T, v string)

type DynamicValues map[string]DynamicValue

func ValueEqual(val string) DynamicValue {
	return func(t *testing.T, v string) {
		assert.Equal(t, val, v)
	}
}

func ValueNotEqual(val string) DynamicValue {
	return func(t *testing.T, v string) {
		assert.NotEqual(t, val, v)
	}
}

func ValueUUID() DynamicValue {
	return func(t *testing.T, v string) {
		_, err := uuid.Parse(v)
		assert.NoError(t, err)
	}
}

func ValueTime() DynamicValue {
	return func(t *testing.T, v string) {
		_, err := time.Parse(time.RFC3339, v)
		assert.NoError(t, err)
	}
}

func ValueTimeInPastDuration(dur time.Duration) DynamicValue {
	return func(t *testing.T, v string) {
		ti, err := time.Parse(time.RFC3339, v)
		assert.NoError(t, err)
		assert.WithinRange(t, ti, time.Now().Add(-dur), time.Now())
	}
}

func ValueRegexp(rx any) DynamicValue {
	return func(t *testing.T, v string) {
		assert.Regexp(t, rx, v)
	}
}
