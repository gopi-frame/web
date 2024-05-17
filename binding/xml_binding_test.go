package binding

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXML(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var container = &struct {
			XMLName xml.Name `xml:"root"`
			Name    string   `xml:"name"`
			Address string   `xml:"address"`
			Age     int      `xml:"age"`
			Valid   bool     `xml:"valid"`
			Tags    []string `xml:"tags"`
		}{}

		assert.Nil(t, XML(r, container))
		assert.Equal(t, "wardonne", container.Name)
		assert.Equal(t, "shanghai", container.Address)
		assert.Equal(t, 10, container.Age)
		assert.Equal(t, true, container.Valid)
		assert.Equal(t, []string{"a", "b"}, container.Tags)
	}))
	defer ts.Close()
	r := strings.NewReader(`<root>
		<name>wardonne</name>
		<address>shanghai</address>
		<age>10</age>
		<valid>true</valid>
		<tags>a</tags>
		<tags>b</tags>
	</root>`)
	resp, err := http.Post(ts.URL, "application/json", r)
	assert.Nil(t, err)
	defer resp.Body.Close()
}
