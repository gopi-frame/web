package response

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderResponse(t *testing.T) {
	reader := strings.NewReader("helloworld")
	response := NewResponse(200).Reader(reader)
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
	assert.Contains(t, result.Header.Get("content-type"), "text/plain")
	assert.Equal(t, "helloworld", string(content))
}
