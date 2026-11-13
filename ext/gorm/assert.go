package golfgorm

import (
	"encoding/json"
	"slices"

	"github.com/gavidroselj/golf"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	modifierDatabase = "db"
)

func AssertTable(t golf.TestingT, db *gorm.DB, model any, dynamicValues golf.DynamicValues, modifiers ...string) {
	result := db.Order(clause.OrderByColumn{Column: clause.PrimaryColumn}).Find(&model)
	assert.NoError(t, result.Error)

	data, err := json.Marshal(model)
	assert.NoError(t, err)

	modifierTableName := result.Statement.Schema.Table

	modifiers = slices.Insert(modifiers, 0, modifierDatabase, modifierTableName)

	golf.AssertJson(t, data, dynamicValues, modifiers...)
}
