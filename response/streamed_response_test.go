package response

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamedResponse(t *testing.T) {
	var i int
	response := New(200).Stream(func(w io.Writer) bool {
		w.Write([]byte(fmt.Sprint(i)))
		i++
		return i < 10
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	response.Send(recorder, request)
	result := recorder.Result()
	body := result.Body
	content, err := io.ReadAll(body)
	assert.Nil(t, err)
	defer func() { body.Close() }()
	assert.Equal(t, "0123456789", string(content))
}
