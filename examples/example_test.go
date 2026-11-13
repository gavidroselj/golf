package examples

import (
	"crypto/rand"
	"testing"
	"time"
	"uuid"

	"github.com/gavidroselj/golf"
)

type ExampleStruct struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func TestExample(t *testing.T) {
	// Assert a string matches a golden file
	golf.AssertString(t, "abc", "string")

	// Assert bytes matches a golden file
	golf.AssertBytes(t, []byte{'d', 'e', 'f'}, "bytes")

	// Marshall data to JSON and assert against a golden file
	data := ExampleStruct{
		Field1: "abc",
		Field2: "def",
	}
	golf.AssertJson(t, data, nil, "json")
}

type DynamicStruct struct {
	Random string    `json:"random"`
	Time   time.Time `json:"time"`
	UUID   uuid.UUID `json:"uuid"`
}

func TestDynamic(t *testing.T) {
	randString := rand.Text()

	data := DynamicStruct{
		Random: randString,
		Time:   time.Now().Add(-200 * time.Millisecond),
		UUID:   uuid.New(),
	}

	golf.AssertJson(t, data, golf.DynamicValues{
		"$.random": golf.ValueEqual(randString),
		"$.time":   golf.ValueTimeInPastDuration(3 * time.Second),
		"$.uuid":   golf.ValueUUID(),
	}, "json")
}
