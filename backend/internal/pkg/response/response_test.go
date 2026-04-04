package response

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseStruct(t *testing.T) {
	resp := Response{
		Code: 200,
		Msg:  "success",
		Data: map[string]string{"key": "value"},
	}

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Msg)
	assert.NotNil(t, resp.Data)
}

func TestResponseJSONSerialization(t *testing.T) {
	resp := Response{
		Code: 200,
		Msg:  "success",
		Data: "test data",
	}

	jsonBytes, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded Response
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, resp.Code, decoded.Code)
	assert.Equal(t, resp.Msg, decoded.Msg)
}

func TestResponseWithNilData(t *testing.T) {
	resp := Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	}

	jsonBytes, err := json.Marshal(resp)
	assert.NoError(t, err)

	// Data should be omitted when nil due to omitempty
	assert.NotContains(t, string(jsonBytes), `"data":null`)
}

func TestResponseErrorCode(t *testing.T) {
	resp := Response{
		Code: 400,
		Msg:  "bad request",
	}

	jsonBytes, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded Response
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, 400, decoded.Code)
	assert.Equal(t, "bad request", decoded.Msg)
}

func TestResponseWithComplexData(t *testing.T) {
	complexData := map[string]interface{}{
		"users": []map[string]string{
			{"id": "1", "name": "Alice"},
			{"id": "2", "name": "Bob"},
		},
		"total": 2,
	}

	resp := Response{
		Code: 200,
		Msg:  "success",
		Data: complexData,
	}

	jsonBytes, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded Response
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), "users")
	assert.Contains(t, string(jsonBytes), "Alice")
}
