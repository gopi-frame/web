package binding

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForm(t *testing.T) {
	t.Run("PostForm", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var container = &struct {
				Name    string   `form:"name"`
				Address string   `form:"address"`
				Age     int      `form:"age"`
				Valid   bool     `form:"valid"`
				Tags    []string `form:"tags"`
			}{}
			assert.Nil(t, Form(r, container))
			assert.Equal(t, "wardonne", container.Name)
			assert.Equal(t, "shanghai", container.Address)
			assert.Equal(t, 10, container.Age)
			assert.Equal(t, true, container.Valid)
			assert.Equal(t, []string{"a", "b"}, container.Tags)
		}))
		defer ts.Close()
		values := make(url.Values)
		values.Add("name", "wardonne")
		values.Add("address", "shanghai")
		values.Add("age", "10")
		values.Add("valid", "true")
		values.Add("tags", "a")
		values.Add("tags", "b")
		resp, err := http.PostForm(ts.URL, values)
		assert.Nil(t, err)
		defer resp.Body.Close()
	})

	t.Run("MultipartForm-OnlyValues", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var container = &struct {
				Name    string   `form:"name"`
				Address string   `form:"address"`
				Age     int      `form:"age"`
				Valid   bool     `form:"valid"`
				Tags    []string `form:"tags"`
			}{}
			assert.Nil(t, Form(r, container))
			assert.Equal(t, "wardonne", container.Name)
			assert.Equal(t, "shanghai", container.Address)
			assert.Equal(t, 10, container.Age)
			assert.Equal(t, true, container.Valid)
			assert.Equal(t, []string{"a", "b"}, container.Tags)
		}))
		defer ts.Close()
		b := bytes.NewBuffer([]byte{})
		w := multipart.NewWriter(b)
		assert.Nil(t, w.WriteField("name", "wardonne"))
		assert.Nil(t, w.WriteField("address", "shanghai"))
		assert.Nil(t, w.WriteField("age", "10"))
		assert.Nil(t, w.WriteField("valid", "true"))
		assert.Nil(t, w.WriteField("tags", "a"))
		assert.Nil(t, w.WriteField("tags", "b"))
		w.Close()
		resp, err := http.Post(ts.URL, w.FormDataContentType(), b)
		assert.Nil(t, err)
		defer resp.Body.Close()
	})

	t.Run("MultipartForm-OnlyFiles", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			container := &struct {
				File1 *UploadedFile `form:"file1"`
				File2 *UploadedFile `form:"file2"`
			}{}
			assert.Nil(t, Form(r, container))
			content1, err := container.File1.Content()
			assert.Nil(t, err)
			assert.Equal(t, "filename1.txt", container.File1.Name())
			assert.Equal(t, "hello world in file1", string(content1))

			content2, err := container.File2.Content()
			assert.Nil(t, err)
			assert.Equal(t, "filename2.txt", container.File2.Name())
			assert.Equal(t, "hello world in file2", string(content2))
		}))
		defer ts.Close()
		b := bytes.NewBuffer([]byte{})
		w := multipart.NewWriter(b)

		p1, err := w.CreateFormFile("file1", "filename1.txt")
		assert.Nil(t, err)
		r1 := bytes.NewReader([]byte("hello world in file1"))
		_, err = r1.Seek(0, 0)
		assert.Nil(t, err)
		_, err = io.Copy(p1, r1)
		assert.Nil(t, err)

		p2, err := w.CreateFormFile("file2", "filename2.txt")
		assert.Nil(t, err)
		r2 := bytes.NewReader([]byte("hello world in file2"))
		_, err = r2.Seek(0, 0)
		assert.Nil(t, err)
		_, err = io.Copy(p2, r2)
		assert.Nil(t, err)

		w.Close()
		resp, err := http.Post(ts.URL, w.FormDataContentType(), b)
		assert.Nil(t, err)
		defer resp.Body.Close()
	})

	t.Run("MultipartForm-FileSlice", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			container := &struct {
				Files UploadedFiles `form:"file[]"`
			}{}

			assert.Nil(t, Form(r, container))
			content1, err := container.Files.Get(0).Content()
			assert.Nil(t, err)
			assert.Equal(t, "filename1.txt", container.Files.Get(0).Name())
			assert.Equal(t, "hello world in file1", string(content1))

			content2, err := container.Files.Get(1).Content()
			assert.Nil(t, err)
			assert.Equal(t, "filename2.txt", container.Files.Get(1).Name())
			assert.Equal(t, "hello world in file2", string(content2))
		}))
		defer ts.Close()
		b := bytes.NewBuffer([]byte{})
		w := multipart.NewWriter(b)

		p1, err := w.CreateFormFile("file[]", "filename1.txt")
		assert.Nil(t, err)
		r1 := bytes.NewReader([]byte("hello world in file1"))
		_, err = r1.Seek(0, 0)
		assert.Nil(t, err)
		_, err = io.Copy(p1, r1)
		assert.Nil(t, err)

		p2, err := w.CreateFormFile("file[]", "filename2.txt")
		assert.Nil(t, err)
		r2 := bytes.NewReader([]byte("hello world in file2"))
		_, err = r2.Seek(0, 0)
		assert.Nil(t, err)
		_, err = io.Copy(p2, r2)
		assert.Nil(t, err)

		w.Close()
		resp, err := http.Post(ts.URL, w.FormDataContentType(), b)
		assert.Nil(t, err)
		defer resp.Body.Close()
	})

	t.Run("MultipartForm-ValueAndFile", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			container := &struct {
				File1   *UploadedFile `form:"file1"`
				File2   *UploadedFile `form:"file2"`
				Name    string        `form:"name"`
				Address string        `form:"address"`
				Age     int           `form:"age"`
				Valid   bool          `form:"valid"`
				Tags    []string      `form:"tags"`
			}{}

			assert.Nil(t, Form(r, container))
			content1, err := container.File1.Content()
			assert.Nil(t, err)
			assert.Equal(t, "filename1.txt", container.File1.Name())
			assert.Equal(t, "hello world in file1", string(content1))

			content2, err := container.File2.Content()
			assert.Nil(t, err)
			assert.Equal(t, "filename2.txt", container.File2.Name())
			assert.Equal(t, "hello world in file2", string(content2))

			assert.Nil(t, Form(r, container))
			assert.Equal(t, "wardonne", container.Name)
			assert.Equal(t, "shanghai", container.Address)
			assert.Equal(t, 10, container.Age)
			assert.Equal(t, true, container.Valid)
			assert.Equal(t, []string{"a", "b"}, container.Tags)
		}))

		defer ts.Close()
		b := bytes.NewBuffer([]byte{})
		w := multipart.NewWriter(b)
		assert.Nil(t, w.WriteField("name", "wardonne"))
		assert.Nil(t, w.WriteField("address", "shanghai"))
		assert.Nil(t, w.WriteField("age", "10"))
		assert.Nil(t, w.WriteField("valid", "true"))
		assert.Nil(t, w.WriteField("tags", "a"))
		assert.Nil(t, w.WriteField("tags", "b"))

		p1, err := w.CreateFormFile("file1", "filename1.txt")
		assert.Nil(t, err)
		r1 := bytes.NewReader([]byte("hello world in file1"))
		_, err = r1.Seek(0, 0)
		assert.Nil(t, err)
		_, err = io.Copy(p1, r1)
		assert.Nil(t, err)

		p2, err := w.CreateFormFile("file2", "filename2.txt")
		assert.Nil(t, err)
		r2 := bytes.NewReader([]byte("hello world in file2"))
		_, err = r2.Seek(0, 0)
		assert.Nil(t, err)
		_, err = io.Copy(p2, r2)
		assert.Nil(t, err)
		w.Close()
		resp, err := http.Post(ts.URL, w.FormDataContentType(), b)
		assert.Nil(t, err)
		defer resp.Body.Close()
	})
}
