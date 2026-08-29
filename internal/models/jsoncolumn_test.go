package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONColumnValue(t *testing.T) {
	data := map[string]any{"foo": "bar"}
	jc := JSONColumn[map[string]any]{Val: data}

	val, err := jc.Value()
	assert.NoError(t, err)

	expected, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, val)
}

func TestJSONColumnScan(t *testing.T) {
	data := `{"foo":"bar"}`
	var jc JSONColumn[map[string]any]

	err := jc.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, "bar", jc.Val["foo"])

	err = jc.Scan([]byte(data))
	assert.NoError(t, err)
	assert.Equal(t, "bar", jc.Val["foo"])

	err = jc.Scan(nil)
	assert.Error(t, err)

	err = jc.Scan(123)
	assert.Error(t, err)
}
