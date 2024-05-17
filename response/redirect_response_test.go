package response

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirectResponse(t *testing.T) {
	response := NewResponse(302).Redirect("https://target.com")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	response.Send(recorder, request)
	result := recorder.Result()
	assert.Equal(t, 302, result.StatusCode)
	assert.Equal(t, "https://target.com", result.Header.Get("location"))
}
