package response

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileResponse(t *testing.T) {
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test_response")
	assert.Nil(t, err)
	tempFile.Write([]byte("helloworld"))
	response := NewResponse(200).File(tempFile.Name())
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	response.Send(recorder, request)
	result := recorder.Result()
	body := result.Body
	content, err := io.ReadAll(body)
	assert.Nil(t, err)
	defer func() { body.Close() }()
	assert.Equal(t, 200, result.StatusCode)
	assert.Contains(t, result.Header.Get("content-type"), "text/plain")
	assert.Equal(t, "helloworld", string(content))
}
