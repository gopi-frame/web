package response

import (
	"net/http"

	"github.com/gopi-frame/codec"
	"github.com/gopi-frame/web/mimetype"
)

// TOMLResponse used to sends TOML-encoded data
type TOMLResponse struct {
	*Response
	data any
}

// SetContent sets response body content
func (tomlResponse *TOMLResponse) SetContent(data any) {
	tomlResponse.data = data
}

// Send sends the response
func (tomlResponse *TOMLResponse) Send(w http.ResponseWriter, r *http.Request) {
	tomlBytes, err := codec.Marshal(codec.TOML, tomlResponse.data)
	if err != nil {
		panic(err)
	}
	tomlResponse.content = tomlBytes
	tomlResponse.SetHeader("Content-Type", mimetype.TOML)
	tomlResponse.Response.Send(w, r)
}
