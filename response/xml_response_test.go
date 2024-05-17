package response

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLResponse(t *testing.T) {
	type data struct {
		Key1 string `xml:"key1"`
		Key2 string `xml:"Key2"`
	}
	data1 := data{
		Key1: "value1",
		Key2: "value2",
	}
	response := NewResponse(200).XML(data1)
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
	assert.Equal(t, "application/xml", result.Header.Get("content-type"))
	assert.Equal(t, `<data><key1>value1</key1><Key2>value2</Key2></data>`, string(content))
}
