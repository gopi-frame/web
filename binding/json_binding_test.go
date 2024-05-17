package binding

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	t.Run("ToStruct", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var container = &struct {
				Name    string   `json:"name"`
				Address string   `json:"address"`
				Age     int      `json:"age"`
				Valid   bool     `json:"valid"`
				Tags    []string `json:"tags"`
			}{}

			assert.Nil(t, JSON(r, container))
			assert.Equal(t, "wardonne", container.Name)
			assert.Equal(t, "shanghai", container.Address)
			assert.Equal(t, 10, container.Age)
			assert.Equal(t, true, container.Valid)
			assert.Equal(t, []string{"a", "b"}, container.Tags)
		}))
		defer ts.Close()
		r := strings.NewReader(`{"name":"wardonne", "address":"shanghai", "age":10, "valid":true, "tags": ["a", "b"]}`)
		resp, err := http.Post(ts.URL, "application/json", r)
		assert.Nil(t, err)
		defer resp.Body.Close()
	})

	t.Run("ToMap", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var container = map[string]any{}

			assert.Nil(t, JSON(r, &container))
			assert.Equal(t, "wardonne", container["name"])
			assert.Equal(t, "shanghai", container["address"])
			assert.EqualValues(t, 10, container["age"])
			assert.Equal(t, true, container["valid"])
			assert.ElementsMatch(t, []string{"a", "b"}, container["tags"])
		}))
		defer ts.Close()
		r := strings.NewReader(`{"name":"wardonne", "address":"shanghai", "age":10, "valid":true, "tags": ["a", "b"]}`)
		resp, err := http.Post(ts.URL, "application/json", r)
		assert.Nil(t, err)
		defer resp.Body.Close()
	})
}
