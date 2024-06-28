package response

import (
	"encoding/xml"
	"net/http"

	"github.com/gopi-frame/web/mimetype"
)

// XMLResponse used to send XML-encoded data
type XMLResponse struct {
	*Response
	data any
}

// SetContent sets response body content
func (xmlResponse *XMLResponse) SetContent(data any) {
	xmlResponse.data = data
}

// Send sends the response
func (xmlResponse *XMLResponse) Send(w http.ResponseWriter, r *http.Request) {
	xmlBytes, err := xml.Marshal(xmlResponse.data)
	if err != nil {
		panic(err)
	}
	xmlResponse.content = xmlBytes
	xmlResponse.SetHeader("Content-Type", mimetype.XML)
	xmlResponse.Response.Send(w, r)
}
