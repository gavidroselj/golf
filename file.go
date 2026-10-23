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

func getFilePath(t TestingT, modifier *Modifier) string {
	if modifier == nil {
		return fmt.Sprintf("%s/%s.golden", config.BasePath, t.Name())
	}

	return fmt.Sprintf("%s/%s.%s.golden", config.BasePath, t.Name(), modifier)
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
	if !config.inhibitUpdate && *updateGoldenFiles {
		writeFile(t, filePath, data)
	}
}

type Modifier struct {
	modifiers []string
}

func (m *Modifier) String() string {
	return strings.Join(m.modifiers, ".")
}

func Modifiers(modifiers ...string) *Modifier {
	return &Modifier{
		modifiers: modifiers,
	}
}
