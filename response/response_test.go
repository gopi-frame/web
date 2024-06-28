package response

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	response := New(200, "hello")
	assert.Equal(t, 200, response.StatusCode())
	assert.EqualValues(t, "hello", response.Content())
	response.SetStatusCode(201)
	assert.Equal(t, 201, response.StatusCode())
	response.SetContent("helloworld")
	assert.Equal(t, "helloworld", response.Content())
	response.SetHeader("content-type", "application/json")
	assert.True(t, response.HasHeader("content-type"))
	assert.Equal(t, "application/json", response.Header("content-type"))
	expectHeader := make(http.Header)
	expectHeader.Set("content-type", "application/json")
	assert.EqualValues(t,
		expectHeader,
		response.Headers())
	response.SetHeaders(map[string]string{"content-type": "application/xml"})
	assert.True(t, response.HasHeader("content-type"))
	assert.Equal(t, "application/xml", response.Header("content-type"))
	expectHeader.Set("content-type", "application/xml")
	assert.EqualValues(t,
		expectHeader,
		response.Headers())
	cookie := new(http.Cookie)
	cookie.Name = "language"
	cookie.Value = "zh-CN"
	response.SetCookie(cookie)
	assert.Equal(t, []*http.Cookie{cookie}, response.Cookies())
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	response.Send(recorder, request)
	assert.Equal(t, 201, recorder.Result().StatusCode)
	assert.Equal(t, 1, len(recorder.Result().Cookies()))
	assert.EqualValues(t, cookie.Name, recorder.Result().Cookies()[0].Name)
	assert.EqualValues(t, cookie.Value, recorder.Result().Cookies()[0].Value)
	assert.Equal(t, "application/xml", recorder.Header().Get("content-type"))
	body := recorder.Result().Body
	data, err := io.ReadAll(body)
	assert.Nil(t, err)
	defer func() {
		body.Close()
	}()
	assert.Equal(t, "helloworld", string(data))
}
