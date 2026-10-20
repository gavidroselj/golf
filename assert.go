package golf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertRaw(t *testing.T, data string) {
	filePath := getFilePath(t)

	updateFileIfFlagSet(t, filePath, []byte(data))

	goldenData := readFile(t, filePath)
	assert.Equal(t, goldenData, data)
}
