package utils_test

import (
	"encoding/json"
	"quantito-engine/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string
	Value int
}

func TestJsonUnmarshalScalar(t *testing.T) {
	jsonBytes, _ := json.Marshal(5)
	out, err := utils.UnmarshalJson[int](jsonBytes)
	assert.NoError(t, err)
	assert.Equal(t, 5, out)
}

func TestJsonUnmarshalStruct(t *testing.T) {
	source := &TestStruct{
		Name:  "A",
		Value: 1,
	}

	jsonBytes, _ := json.Marshal(source)
	out, err := utils.UnmarshalJson[*TestStruct](jsonBytes)
	assert.NoError(t, err)
	assert.Equal(t, source, out)
}

func TestJsonUnmarshalSlice(t *testing.T) {
	source := []*TestStruct{{
		Name:  "A",
		Value: 1,
	}, {
		Name:  "B",
		Value: 2,
	}}

	jsonBytes, _ := json.Marshal(source)
	out, err := utils.UnmarshalJson[[]*TestStruct](jsonBytes)
	assert.NoError(t, err)
	assert.Equal(t, source, out)
}

func TestJsonUnmarshalMap(t *testing.T) {
	source := map[string]*TestStruct{
		"key": {
			Name:  "A",
			Value: 1,
		},
	}

	jsonBytes, _ := json.Marshal(source)
	out, err := utils.UnmarshalJson[map[string]*TestStruct](jsonBytes)
	assert.NoError(t, err)
	assert.Equal(t, source, out)
}

func TestJsonToStructWithTags(t *testing.T) {
	type TaggedStruct struct {
		Name     string `json:"name"`
		Value    int    `json:"value"`
		Optional any    `json:"optional,omitempty"`
		Private  any    `json:"-"`
	}

	expectedOutput := &TaggedStruct{
		Name:     "Test",
		Value:    42,
		Optional: nil,
		Private:  nil,
	}

	jsonBytes := []byte(`{"name":"Test","value":42}`)
	out, err := utils.UnmarshalJson[*TaggedStruct](jsonBytes)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, out)

}

func TestJsonUnmarshalToInvalidType(t *testing.T) {
	source := &TestStruct{
		Name:  "A",
		Value: 1,
	}

	jsonBytes, _ := json.Marshal(source)
	_, err := utils.UnmarshalJson[int](jsonBytes)
	assert.Error(t, err)
}

func TestJsonUnmarshalError(t *testing.T) {
	invalidJsonBytes := []byte(`{"Name`)
	_, err := utils.UnmarshalJson[*TestStruct](invalidJsonBytes)
	assert.Error(t, err)
}
