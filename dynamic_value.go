package golf

import (
	"bytes"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// DynamicValue is a validation function for a value, that will change between different test runs.
// The function validates that the existing value satisfies certain criteria.
// The Dynamic Value will then be replaced by a configured placeholder, which will be written to the golden file.
type DynamicValue func(t TestingT, v string)

// DynamicValues is a collection of multiple DynamicValue objects keyed by the JSONPath where they appear.
type DynamicValues map[string]DynamicValue

// ValueEqual checks whether the value is equal to `val`.
func ValueEqual(val string) DynamicValue {
	return func(t TestingT, v string) {
		assert.Equal(t, val, v)
	}
}

// ValueNotEqual checks whether the value is not equal to `val`.
func ValueNotEqual(val string) DynamicValue {
	return func(t TestingT, v string) {
		assert.NotEqual(t, val, v)
	}
}

// ValueUUID checks whether the value is a valid UUID.
func ValueUUID() DynamicValue {
	return func(t TestingT, v string) {
		_, err := uuid.Parse(v)
		assert.NoError(t, err)
	}
}

// ValueTime checks whether the value is a valid RFC3339 timestamp.
func ValueTime() DynamicValue {
	return func(t TestingT, v string) {
		_, err := time.Parse(time.RFC3339, v)
		assert.NoError(t, err)
	}
}

// ValueTimeInPastDuration checks whether the value is a valid RFC3339 timestamp between Now and `dur` before now.
// Times in the future up to 100 milliseconds are accepted to prevent flakey tests.
func ValueTimeInPastDuration(dur time.Duration) DynamicValue {
	return func(t TestingT, v string) {
		ti, err := time.Parse(time.RFC3339, v)
		assert.NoError(t, err)

		now := time.Now()
		assert.WithinRange(t, ti, now.Add(-dur), now.Add(100*time.Millisecond))
	}
}

// ValueRegexp checks whether the value satisfies the given expression.
func ValueRegexp(rx any) DynamicValue {
	return func(t TestingT, v string) {
		assert.Regexp(t, rx, v)
	}
}

// ValueBcryptPassword checks whether the value is a valid bcrypt hash for the provided `password`.
func ValueBcryptPassword(password string) DynamicValue {
	return func(t TestingT, v string) {
		storedPassword := bytes.Trim([]byte(v), "\"")
		err := bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
		assert.NoError(t, err)
	}
}
