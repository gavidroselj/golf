package golf

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AssertRawWithModifiers(t *testing.T, data string, modifier *Modifier) {
	filePath := getFilePath(t, modifier)

	updateFileIfFlagSet(t, filePath, []byte(data))

	goldenData := readFile(t, filePath)
	assert.Equal(t, goldenData, data)
}

func AssertJSON(t *testing.T, data any, modifier *Modifier, dynamicValues DynamicValues) {
	jsonData, err := json.Marshal(data, json.Deterministic(true), jsontext.Multiline(true))
	assert.NoError(t, err, "Failed to marshall data to json", data)

	for path, dynamicValue := range dynamicValues {
		parsedPath, err := jsonparser.ParsePath(path)
		assert.NoError(t, err, "Failed to parse JSON path", path)

		value, err := jsonparser.GetString(jsonData, parsedPath...)
		if err != nil {
			// Path does not exist, skip
			continue
		}

		dynamicValue(t, value)

		jsonData, err = jsonparser.SetString(jsonData, config.DynamicValuePlaceholder, parsedPath...)
		assert.NoError(t, err, "Failed to replace dynamic value with placeholder")
	}

	filePath := getFilePath(t, modifier)

	updateFileIfFlagSet(t, filePath, jsonData)

	goldenData := readFile(t, filePath)
	assert.JSONEq(t, string(goldenData), string(jsonData))
}

func AssertGormTable(t *testing.T, db *gorm.DB, model any, dynamicValues DynamicValues) {
	result := db.Order(clause.OrderByColumn{Column: clause.PrimaryColumn}).Find(&model)
	assert.NoError(t, result.Error)

	data, err := json.Marshal(model)
	assert.NoError(t, err)

	modifier := Modifiers("db", result.Statement.Schema.Table)

	AssertJSON(t, data, modifier, dynamicValues)
}
