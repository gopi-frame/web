package response

import (
	"net/http"

	"github.com/gopi-frame/codec"
	"github.com/gopi-frame/web/mimetype"
)

// YAMLResponse used to send YAML-encoded response
type YAMLResponse struct {
	*Response
	data any
}

// SetContent sets response body content
func (yamlResponse *YAMLResponse) SetContent(data any) {
	yamlResponse.data = data
}

// Send sends the response
func (yamlResponse *YAMLResponse) Send(w http.ResponseWriter, r *http.Request) {
	yamlBytes, err := codec.Marshal(codec.YAML, yamlResponse.data)
	if err != nil {
		panic(err)
	}
	yamlResponse.content = yamlBytes
	yamlResponse.SetHeader("Content-Type", mimetype.YAML)
	yamlResponse.Response.Send(w, r)
}
