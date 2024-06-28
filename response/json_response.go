package response

import (
	"encoding/json"
	"net/http"

	"github.com/gopi-frame/web/mimetype"
)

// JSONResponse used to response json format data
type JSONResponse struct {
	*Response
	data any
}

// SetContent sets response body content
func (jsonResponse *JSONResponse) SetContent(data any) {
	jsonResponse.data = data
}

// Send sends the response
func (jsonResponse *JSONResponse) Send(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(jsonResponse.data)
	if err != nil {
		panic(err)
	}
	jsonResponse.content = jsonBytes
	jsonResponse.SetHeader("Content-Type", mimetype.JSON)
	jsonResponse.Response.Send(w, r)
}
