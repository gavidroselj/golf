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

func readAndUpdateFile(t *testing.T, data []byte, modifier *Modifier) []byte {
	filePath := getFilePath(t, modifier)
	updateFileIfFlagSet(t, filePath, data)

	return readFile(t, filePath)
}

func AssertBytes(t *testing.T, data []byte, modifier *Modifier) {
	goldenData := readAndUpdateFile(t, data, modifier)
	assert.Equal(t, goldenData, data)
}

func AssertString(t *testing.T, data string, modifier *Modifier) {
	goldenData := string(readAndUpdateFile(t, []byte(data), modifier))
	assert.Equal(t, goldenData, data)
}

func checkDynamicValues(t *testing.T, jsonData []byte, dynamicValues DynamicValues) []byte {
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

	return jsonData
}

func AssertJson(t *testing.T, data any, modifier *Modifier, dynamicValues DynamicValues) {
	jsonData, err := json.Marshal(data, json.Deterministic(true), jsontext.Multiline(true))
	assert.NoError(t, err, "Failed to marshall data to json", data)

	nonDynamicData := checkDynamicValues(t, jsonData, dynamicValues)

	goldenData := readAndUpdateFile(t, nonDynamicData, modifier)
	assert.JSONEq(t, string(goldenData), string(nonDynamicData))
}

func AssertGormTable(t *testing.T, db *gorm.DB, model any, dynamicValues DynamicValues) {
	result := db.Order(clause.OrderByColumn{Column: clause.PrimaryColumn}).Find(&model)
	assert.NoError(t, result.Error)

	data, err := json.Marshal(model)
	assert.NoError(t, err)

	modifier := Modifiers("db", result.Statement.Schema.Table)

	AssertJson(t, data, modifier, dynamicValues)
}
