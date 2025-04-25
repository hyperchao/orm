package tag

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func parseFunc(tag string) (string, string) {
	return tag, tag
}

func TestParser_Parse_Set_Interface(t *testing.T) {

	type InnerModel struct {
		Name       string `test:"name"`
		Age        int    `test:"age"`
		unexported string `test:"unexported"`
	}

	type OuterModel struct {
		*InnerModel
		Score float32 `test:"score"`

		Other *OuterModel // test circular reference
	}

	parser := NewParser(parseFunc)

	modelValues := parser.Parse("test", "not struct")
	assert.Equal(t, 0, modelValues.Len())
	assert.Nil(t, modelValues.Get("not exist"))

	var innerModel InnerModel

	modelValues = parser.Parse("test", &innerModel)
	assert.Equal(t, 2, modelValues.Len())
	assert.Nil(t, modelValues.Get("not exist"))
	assert.Equal(t, "name", modelValues.Get("name").Meta().Name())
	assert.Equal(t, "name", modelValues.Get("name").Meta().Attrs())
	assert.Equal(t, reflect.String, modelValues.Get("name").Meta().Type().Kind())
	for name, value := range modelValues.Iter() {
		assert.Equal(t, name, value.Meta().Name())
	}

	nameVal := modelValues.Get("name")
	nameVal.Set("123")
	assert.Equal(t, "123", innerModel.Name)

	var outerModel OuterModel
	modelValues = parser.Parse("test", &outerModel)
	assert.Equal(t, 3, modelValues.Len())
	modelValues.Get("name").Set("xxx")
	modelValues.Get("age").Set(18)
	modelValues.Get("score").Set(100.0)
	assert.Equal(t, "xxx", outerModel.Name)
	assert.Equal(t, 18, outerModel.Age)
	assert.Equal(t, float32(100), outerModel.Score)

	innerValues := parser.Parse("test", outerModel.InnerModel)
	assert.Equal(t, 2, innerValues.Len())
	assert.Equal(t, "xxx", innerValues.Get("name").Interface())
	innerValues.Get("name").Set("yyy")
	assert.Equal(t, "yyy", modelValues.Get("name").Interface())
}
