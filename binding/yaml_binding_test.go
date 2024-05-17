package binding

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLParser_Parse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var container = &struct {
			Name    string   `yaml:"name"`
			Address string   `yaml:"address"`
			Age     int      `yaml:"age"`
			Valid   bool     `yaml:"valid"`
			Tags    []string `yaml:"tags"`
		}{}
		assert.Nil(t, YAML(r, container))
		assert.Equal(t, "wardonne", container.Name)
		assert.Equal(t, "shanghai", container.Address)
		assert.Equal(t, 10, container.Age)
		assert.Equal(t, true, container.Valid)
		assert.Equal(t, []string{"a", "b"}, container.Tags)
	}))
	defer ts.Close()
	r := strings.NewReader(`name: "wardonne"
address: "shanghai"
age: 10
valid: true
tags:
 - "a"
 - "b"`)
	resp, err := http.Post(ts.URL, "application/json", r)
	assert.Nil(t, err)
	defer resp.Body.Close()
}
