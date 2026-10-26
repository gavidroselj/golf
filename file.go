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

func getFilePath(t TestingT, modifier *Modifiers) string {
	cfg := getConfig(t)

	testPath := t.Name()
	if cfg.IgnoreSubtestNames {
		testPath = topLevelTestName(t)
	}

	if modifier == nil {
		return fmt.Sprintf("%s/%s.golden", cfg.BasePath, testPath)
	}

	return fmt.Sprintf("%s/%s.%s.golden", cfg.BasePath, testPath, modifier)
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

// Modifiers is used to differentiate multiple different golden files used in the same test.
// It appends all modifiers to the file name, ending with the .golden extension.
// Example: A test named "TestSomething" with subtest name "abc" and modifiers ["mod1"] would be saved in BasePath/TestSomething/abc.mod1.golden
type Modifiers struct {
	modifiers []string
}

func (m *Modifiers) String() string {
	return strings.Join(m.modifiers, ".")
}

// NewModifiers is a utility function that creates a Modifiers object from the provided strings.
func NewModifiers(modifiers ...string) *Modifiers {
	return &Modifiers{
		modifiers: modifiers,
	}
}
