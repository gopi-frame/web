package response

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONResponse(t *testing.T) {
	response := New(200).JSON(map[string]any{
		"key1": "value1",
		"key2": "value2",
	})
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
	assert.Equal(t, "application/json", result.Header.Get("content-type"))
	assert.JSONEq(t, `{"key1":"value1","key2":"value2"}`, string(content))
}
