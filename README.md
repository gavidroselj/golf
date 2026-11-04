# golf

Golden file testing framework for Go

## Usage

Install the package via go get.

```bash
go get github.com/gavidroselj/golf
```

Add assertion statements to your tests as needed.

```go
package examples

import (
	"testing"

	"github.com/gavidroselj/golf"
)

type ExampleStruct struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func TestExample(t *testing.T) {
	// Assert a string matches a golden file
	golf.AssertString(t, "abc", golf.NewModifiers("string"))

	// Assert bytes matches a golden file
	golf.AssertBytes(t, []byte{'d', 'e', 'f'}, golf.NewModifiers("bytes"))

	// Marshall data to JSON and assert against a golden file
	data := ExampleStruct{
		Field1: "abc",
		Field2: "def",
	}
	golf.AssertJson(t, data, golf.NewModifiers("json"), nil)
}
```

Run the tests.

```bash
go test ./...
```

Update and generate all golden files with the `-update` test flag.

```bash
go test ./... -update
```

### Modifiers

By default, each test gets it's test file name from t.Name(), which is in the format `TestTopLevel/subtest1/subtest2`. This along with the BasePath is used to create the path where the test golden file will be stored.

See configuration options WithBasePath and WithIgnoreSbutestNames to fine tune this behavior.

Modifiers will be appended to the golden file name, thus allowing you to have multiple different golden files in the same test.

### Dynamic values

JSON and all higher level assertions support field values that are different on every test run, but follow specific rules.

When dynamic value, the field value is extracted, compared to the DynamicValue validator and replaced with a placeholder, which is written into the golden file. If the value fails validation, the test will fail.

The dynamic values are provided as a `map[string]DynamicValue`, where the key represents the [JSONpath](https://en.wikipedia.org/wiki/JSONPath) to the target value.

Example usage:

```go
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

	golf.AssertJson(t, data, golf.NewModifiers("json"), golf.DynamicValues{
		"$.random": golf.ValueEqual(randString),
		"$.time":   golf.ValueTimeInPastDuration(3 * time.Second),
		"$.uuid":   golf.ValueUUID(),
	})
}
```

### Configuring

The package has several options that can be configured globally or on a per-top-level-test basis.

Setting global defaults for your entire project

```go
func TestMain(m *testing.M) {
	// Configure defaults, these will be applied to all tests
	golf.SetDefaults(golf.WithBasePath("alternate_testdata"))

	os.Exit(m.Run())
}
```

Modifying the defaults for a specific test and all it's subtests

```go
func TestCustom(t *testing.T) {
	// This test and all subtests will inherit default values and modify them
	golf.Configure(t, golf.WithBasePath("alternate_testdata_2"))

	golf.AssertString(t, "abc", nil)

	t.Run("subtest", func(t *testing.T) {
	    golf.AssertString(t, "def", nil)
	})
}
```
