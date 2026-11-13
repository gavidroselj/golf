package golf

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"slices"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	modifierDatabase = "db"
)

func AssertBytes(t TestingT, data []byte, modifiers ...string) {
	goldenData := readAndUpdateFile(t, data, modifiers...)
	assert.Equal(t, goldenData, data)
}

func AssertString(t TestingT, data string, modifiers ...string) {
	goldenData := string(readAndUpdateFile(t, []byte(data), modifiers...))
	assert.Equal(t, goldenData, data)
}

func checkDynamicValues(t TestingT, jsonData []byte, dynamicValues DynamicValues) []byte {
	cfg := getConfig(t)

	for path, dynamicValue := range dynamicValues {
		parsedPath, err := jsonparser.ParsePath(path)
		assert.NoError(t, err, "Failed to parse JSON path", path)

		value, err := jsonparser.GetString(jsonData, parsedPath...)
		if err != nil {
			// Path does not exist, skip
			continue
		}

		dynamicValue(t, value)

		jsonData, err = jsonparser.SetString(jsonData, cfg.DynamicValuePlaceholder, parsedPath...)
		assert.NoError(t, err, "Failed to replace dynamic value with placeholder")
	}

	return jsonData
}

func AssertJson(t TestingT, data any, dynamicValues DynamicValues, modifiers ...string) {
	jsonData, err := json.Marshal(data, json.Deterministic(true), jsontext.Multiline(true))
	assert.NoError(t, err, "Failed to marshall data to json", data)

	nonDynamicData := checkDynamicValues(t, jsonData, dynamicValues)

	goldenData := readAndUpdateFile(t, nonDynamicData, modifiers...)
	assert.JSONEq(t, string(goldenData), string(nonDynamicData))
}

func AssertGormTable(t TestingT, db *gorm.DB, model any, dynamicValues DynamicValues, modifiers ...string) {
	result := db.Order(clause.OrderByColumn{Column: clause.PrimaryColumn}).Find(&model)
	assert.NoError(t, result.Error)

	data, err := json.Marshal(model)
	assert.NoError(t, err)

	modifierTableName := result.Statement.Schema.Table

	modifiers = slices.Insert(modifiers, 0, modifierDatabase, modifierTableName)

	AssertJson(t, data, dynamicValues, modifiers...)
}
