package golf

import (
	"testing"
	"time"
	"uuid"
)

const variantTestCount = 10000

func TestValueEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		val := ValueEqual("abc")
		val(t, "abc")
	})

	t.Run("fail", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueEqual("abc")
			val(t, "def")
		})
	})
}

func TestValueNotEqual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		val := ValueNotEqual("abc")
		val(t, "def")
	})

	t.Run("fail", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueNotEqual("abc")
			val(t, "abc")
		})
	})
}

func TestValueUUID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		val := ValueUUID()
		val(t, "7e874f42-cf19-11f1-be6b-bc2411aa2981")
	})

	t.Run("success-random", func(t *testing.T) {
		val := ValueUUID()
		for range variantTestCount {
			id := uuid.NewV4()
			val(t, id.String())
		}
	})

	t.Run("fail-format", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueUUID()
			val(t, "abc")
		})
	})

	t.Run("fail-empty", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueUUID()
			val(t, "")
		})
	})
}

func TestValueTime(t *testing.T) {
	t.Run("success-fixed", func(t *testing.T) {
		val := ValueTime()
		val(t, "1985-04-12T23:20:50.52Z")
	})

	t.Run("success-zone", func(t *testing.T) {
		val := ValueTime()
		val(t, "1996-12-19T16:39:57-08:00")
	})

	t.Run("fail-format", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueTime()
			val(t, "abc")
		})
	})

	t.Run("fail-empty", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueTime()
			val(t, "")
		})
	})
}

func TestValueTimeInPastDuration(t *testing.T) {
	t.Run("success-now", func(t *testing.T) {
		for range variantTestCount {
			val := ValueTimeInPastDuration(3 * time.Second)
			val(t, time.Now().Format(time.RFC3339))
		}
	})

	t.Run("success-past", func(t *testing.T) {
		for range variantTestCount {
			val := ValueTimeInPastDuration(3 * time.Second)
			val(t, time.Now().Add(-time.Millisecond*200).Format(time.RFC3339))
		}
	})

	t.Run("fail-future", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			for range variantTestCount {
				val := ValueTimeInPastDuration(time.Second)
				val(t, time.Now().Add(time.Second*3).Format(time.RFC3339))
			}
		})
	})

	t.Run("fail-format", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueTime()
			val(t, "abc")
		})
	})

	t.Run("fail-empty", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueTime()
			val(t, "")
		})
	})
}

func TestValueRegexp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		val := ValueRegexp("[a-z]+")
		val(t, "abc")
	})

	t.Run("success-submatch", func(t *testing.T) {
		val := ValueRegexp("[a-z]+")
		val(t, "abc def")
	})

	t.Run("fail-no-match", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueRegexp("[a-z]+")
			val(t, "ABC")
		})
	})

	t.Run("fail-empty", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueRegexp("[a-z]+")
			val(t, "")
		})
	})
}

func TestValueBcryptPassword(t *testing.T) {
	// password123
	const hash = "$2y$10$0PBcxOCMl6a0a3UCDW5b2easnSkBuHVP.r3LSwdmtREUWvG73u0Wy"

	t.Run("success", func(t *testing.T) {
		val := ValueBcryptPassword("password123")
		val(t, hash)
	})

	t.Run("fail-wrong-password", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueBcryptPassword("abc")
			val(t, hash)
		})
	})

	t.Run("fail-decrypted", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueBcryptPassword("password123")
			val(t, "password123")
		})
	})

	t.Run("fail-format", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueBcryptPassword("password123")
			val(t, "abc")
		})
	})

	t.Run("fail-empty", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			val := ValueBcryptPassword("password123")
			val(t, "")
		})
	})
}
