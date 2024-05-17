package response

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTOMLResponse(t *testing.T) {
	type data struct {
		Key1 string `toml:"key1"`
		Key2 string `toml:"Key2"`
	}
	data1 := data{
		Key1: "value1",
		Key2: "value2",
	}
	response := NewResponse(200).TOML(data1)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	response.Send(recorder, request)
	result := recorder.Result()
	body := result.Body
	content, err := io.ReadAll(body)
	assert.Nil(t, err)
	defer func() {
		body.Close()
	}()
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "application/toml", result.Header.Get("content-type"))
	assert.Equal(t, "key1 = 'value1'\nKey2 = 'value2'\n", string(content))
}
