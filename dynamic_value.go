package golf

import (
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
)

type DynamicValue func(t TestingT, v string)

type DynamicValues map[string]DynamicValue

func ValueEqual(val string) DynamicValue {
	return func(t TestingT, v string) {
		assert.Equal(t, val, v)
	}
}

func ValueNotEqual(val string) DynamicValue {
	return func(t TestingT, v string) {
		assert.NotEqual(t, val, v)
	}
}

func ValueUUID() DynamicValue {
	return func(t TestingT, v string) {
		_, err := uuid.Parse(v)
		assert.NoError(t, err)
	}
}

func ValueTime() DynamicValue {
	return func(t TestingT, v string) {
		_, err := time.Parse(time.RFC3339, v)
		assert.NoError(t, err)
	}
}

func ValueTimeInPastDuration(dur time.Duration) DynamicValue {
	return func(t TestingT, v string) {
		ti, err := time.Parse(time.RFC3339, v)
		assert.NoError(t, err)

		now := time.Now().Round(time.Second)
		assert.WithinRange(t, ti, now.Add(-dur), now)
	}
}

func ValueRegexp(rx any) DynamicValue {
	return func(t TestingT, v string) {
		assert.Regexp(t, rx, v)
	}
}
