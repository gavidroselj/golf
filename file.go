package golf

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/stretchr/testify/assert"
)

var (
	updateGoldenFiles = flag.Bool("update", false, "update golden files")
)

func getFilePath(t TestingT, modifiers ...string) string {
	cfg := getConfig(t)

	testPath := t.Name()
	if cfg.IgnoreSubtestNames {
		testPath = topLevelTestName(t)
	}

	if len(modifiers) == 0 {
		return fmt.Sprintf("%s/%s.golden", cfg.BasePath, testPath)
	}

	modifierString := strings.Join(modifiers, ".")

	return fmt.Sprintf("%s/%s.%s.golden", cfg.BasePath, testPath, modifierString)
}

func readFile(t TestingT, filePath string) []byte {
	f, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Error loading golden file")
	return f
}

func writeFile(t TestingT, filePath string, data []byte) {
	err := os.MkdirAll(path.Dir(filePath), 0755)
	assert.NoError(t, err, "Failed to create golden file folder")

	err = os.WriteFile(filePath, data, 0644)
	assert.NoError(t, err, "Failed to write golden file")
}

func updateFileIfFlagSet(t TestingT, filePath string, data []byte) {
	cfg := getConfig(t)

	if !cfg.InhibitUpdate && *updateGoldenFiles {
		writeFile(t, filePath, data)
	}
}

func readAndUpdateFile(t TestingT, data []byte, modifiers ...string) []byte {
	filePath := getFilePath(t, modifiers...)
	updateFileIfFlagSet(t, filePath, data)

	return readFile(t, filePath)
}
