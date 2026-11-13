package golf

import (
	"crypto/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Prevent automatic test file updates during testing
	SetDefaults(WithInhibitUpdate(true))

	os.Exit(m.Run())
}

func TestReadAndUpdateFile(t *testing.T) {
	t.Run("match", func(t *testing.T) {
		data := []byte("Raw Data")
		goldenData := readAndUpdateFile(t, []byte(data))

		assert.Equal(t, data, goldenData)
	})
	t.Run("mismatch", func(t *testing.T) {
		data := []byte("Wrong Data")
		goldenData := readAndUpdateFile(t, data)

		assert.NotEqual(t, data, goldenData)
	})
	// TODO: Add file write tests
}

func TestAssertBytes(t *testing.T) {
	Configure(t, WithIgnoreSubtestNames(true))

	t.Run("match", func(t *testing.T) {
		data := []byte("Bytes Data")
		AssertBytes(t, data)
	})
	t.Run("mismatch", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			data := []byte("Wrong Data")
			AssertBytes(t, data)
		})
	})
}

func TestAssertString(t *testing.T) {
	Configure(t, WithIgnoreSubtestNames(true))

	t.Run("match", func(t *testing.T) {
		data := "String Data"
		AssertString(t, data)
	})
	t.Run("mismatch", func(t *testing.T) {
		ExpectFailure(t, func(t TestingT) {
			data := "Wrong Data"
			AssertString(t, data)
		})
	})
}

type TestModel struct {
	StringField string    `json:"string_field"`
	IntField    int       `json:"int_field"`
	BoolField   bool      `json:"bool_field"`
	TimeField   time.Time `json:"time_field"`
}

func newTestModel() TestModel {
	return TestModel{
		StringField: "text",
		IntField:    123,
		BoolField:   true,
		TimeField:   time.Date(2020, 02, 02, 02, 02, 02, 02, time.UTC),
	}
}

func TestAssertJson(t *testing.T) {
	t.Run("static", func(t *testing.T) {
		AssertJson(t, newTestModel(), nil)
	})
	t.Run("dynamic", func(t *testing.T) {
		randString := rand.Text()
		pastTime := time.Now().Add(-time.Second)

		model := newTestModel()
		model.StringField = randString
		model.TimeField = pastTime

		dv := DynamicValues{
			"$.string_field": ValueEqual(randString),
			"$.time_field":   ValueTimeInPastDuration(time.Second * 5),
		}

		AssertJson(t, model, dv)
	})
}
