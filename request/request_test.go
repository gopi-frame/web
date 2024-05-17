package request

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gopi-frame/types"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRequest_Get(t *testing.T) {
	u := "http://test.site/detail?id=1&d[a]=1&d[b]=2"
	r := httptest.NewRequest("GET", u, nil)
	r.Header.Add("custom-header", "test-request")
	req := NewRequest(r, httprouter.Params{
		{Key: "id", Value: "1"},
	})
	req.Set("ck", "cv")
	assert.EqualValues(t, "cv", req.MustGet("ck"))
	assert.Equal(t, "GET", req.Method())
	assert.True(t, req.IsGet())
	assert.False(t, req.IsPost())
	assert.False(t, req.IsPut())
	assert.False(t, req.IsPatch())
	assert.False(t, req.IsDelete())
	assert.False(t, req.IsHead())
	assert.False(t, req.IsConnect())
	assert.False(t, req.IsOptions())
	assert.False(t, req.IsTrace())
	assert.Equal(t, "test.site", req.Host())
	assert.Equal(t, u, req.RequestURI())
	assert.Equal(t, "/detail", req.Path())
	assert.EqualValues(t, "1", req.Query("id"))
	assert.EqualValues(t, "name", req.Query("name", "name"))
	d, ok := req.QueryMap("d")
	assert.True(t, ok)
	assert.EqualValues(t, map[string]types.String{"a": "1", "b": "2"}, d.ToMap())
	assert.EqualValues(t, "1", req.Param("id"))
	assert.EqualValues(t, "test-request", req.Header("custom-header"))
	assert.EqualValues(t, "not-exists", req.Header("not-exists-header", "not-exists"))
	assert.Equal(t, "192.0.2.1", req.ClientIP())
}

func TestRequest_PostForm(t *testing.T) {
	u := "http://test.site/detail?id=1&d[a]=1&d[b]=2"
	form := make(url.Values)
	form.Set("key1", "value1")
	form.Set("d2[1]", "1")
	form.Set("d2[2]", "2")
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", u, body)
	r.Header.Add("custom-header", "test-request")
	r.Header.Add("content-type", "application/x-www-form-urlencoded")
	req := NewRequest(r, httprouter.Params{
		{Key: "id", Value: "1"},
	})
	assert.Equal(t, "POST", req.Method())
	assert.False(t, req.IsGet())
	assert.True(t, req.IsPost())
	assert.False(t, req.IsPut())
	assert.False(t, req.IsPatch())
	assert.False(t, req.IsDelete())
	assert.False(t, req.IsHead())
	assert.False(t, req.IsConnect())
	assert.False(t, req.IsOptions())
	assert.False(t, req.IsTrace())
	assert.Equal(t, "test.site", req.Host())
	assert.Equal(t, u, req.RequestURI())
	assert.Equal(t, "/detail", req.Path())
	assert.EqualValues(t, "1", req.Query("id"))
	assert.EqualValues(t, "name", req.Query("name", "name"))
	d, ok := req.QueryMap("d")
	assert.True(t, ok)
	assert.EqualValues(t, map[string]types.String{"a": "1", "b": "2"}, d.ToMap())
	assert.EqualValues(t, "1", req.Param("id"))
	assert.EqualValues(t, "test-request", req.Header("custom-header"))
	assert.EqualValues(t, "not-exists", req.Header("not-exists-header", "not-exists"))
	d2, ok := req.PostFormMap("d2")
	assert.True(t, ok)
	assert.EqualValues(t, map[string]types.String{"1": "1", "2": "2"}, d2.ToMap())
	assert.EqualValues(t, "value1", req.PostForm("key1"))
	assert.EqualValues(t, "not-exists-value", req.PostForm("not-exists", "not-exists-value"))
	assert.Equal(t, "192.0.2.1", req.ClientIP())
}

func TestRequest_FormData(t *testing.T) {
	u := "http://test.site/detail?id=1&d[a]=1&d[b]=2"
	b := bytes.NewBufferString("")
	w := multipart.NewWriter(b)
	assert.Nil(t, w.WriteField("key1", "value1"))
	p1, err := w.CreateFormFile("file1", "filename1.txt")
	assert.Nil(t, err)
	r1 := bytes.NewReader([]byte("hello world in file1"))
	_, err = r1.Seek(0, 0)
	assert.Nil(t, err)
	_, err = io.Copy(p1, r1)
	assert.Nil(t, err)

	p2, err := w.CreateFormFile("files[]", "files_0.txt")
	assert.Nil(t, err)
	r2 := bytes.NewReader([]byte("hello world in files_0"))
	_, err = r2.Seek(0, 0)
	assert.Nil(t, err)
	_, err = io.Copy(p2, r2)
	assert.Nil(t, err)

	w.Close()
	r := httptest.NewRequest("POST", u, b)
	r.Header.Add("custom-header", "test-request")
	r.Header.Add("content-type", w.FormDataContentType())
	req := NewRequest(r, httprouter.Params{
		{Key: "id", Value: "1"},
	})
	assert.Equal(t, "POST", req.Method())
	assert.False(t, req.IsGet())
	assert.True(t, req.IsPost())
	assert.False(t, req.IsPut())
	assert.False(t, req.IsPatch())
	assert.False(t, req.IsDelete())
	assert.False(t, req.IsHead())
	assert.False(t, req.IsConnect())
	assert.False(t, req.IsOptions())
	assert.False(t, req.IsTrace())
	assert.Equal(t, "test.site", req.Host())
	assert.Equal(t, u, req.RequestURI())
	assert.Equal(t, "/detail", req.Path())
	assert.EqualValues(t, "1", req.Query("id"))
	assert.EqualValues(t, "name", req.Query("name", "name"))
	d, ok := req.QueryMap("d")
	assert.True(t, ok)
	assert.EqualValues(t, map[string]types.String{"a": "1", "b": "2"}, d.ToMap())
	assert.EqualValues(t, "1", req.Param("id"))
	assert.EqualValues(t, "test-request", req.Header("custom-header"))
	assert.EqualValues(t, "not-exists", req.Header("not-exists-header", "not-exists"))
	assert.Equal(t, "192.0.2.1", req.ClientIP())
	f, err := req.File("file1")
	assert.Nil(t, err)
	content, err := f.Content()
	assert.Nil(t, err)
	assert.EqualValues(t, "hello world in file1", string(content))
	fs := req.Files("files[]").ToArray()
	assert.Equal(t, 1, len(fs))
	content, err = fs[0].Content()
	assert.Equal(t, "hello world in files_0", string(content))
	assert.Nil(t, err)
}
