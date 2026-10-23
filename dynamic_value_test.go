package golf

import (
	"testing"
)

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
