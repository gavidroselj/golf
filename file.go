package golf

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	updateGoldenFiles = flag.Bool("update", false, "update golden files")
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func getFilePath(t *testing.T, modifiers ...string) string {
	if len(modifiers) == 0 {
		return fmt.Sprintf("%s/%s.golden", config.BasePath, t.Name())
	}

	modifier := strings.Join(modifiers, ".")
	return fmt.Sprintf("%s/%s.%s.golden", config.BasePath, t.Name(), modifier)
}

func readFile(t *testing.T, filePath string) []byte {
	f, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Error loading golden file")
	return f
}

func writeFile(t *testing.T, filePath string, data []byte) {
	err := os.MkdirAll(path.Dir(filePath), 0755)
	require.NoError(t, err, "Failed to create golden file folder")

	err = os.WriteFile(filePath, data, 0644)
	require.NoError(t, err, "Failed to write golden file")
}

func updateFileIfFlagSet(t *testing.T, filePath string, data []byte) {
	if *updateGoldenFiles {
		writeFile(t, filePath, data)
	}
}
