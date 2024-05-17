package response

import (
	"net/http"

	"github.com/gopi-frame/web/mimetype"
	"google.golang.org/protobuf/proto"
)

// ProtobufResponse used to send a protobuf response
type ProtobufResponse struct {
	*Response
	data any
}

// SetContent sets response body content
//
// NOTICE: data should be [proto.Message]
func (protobufResponse *ProtobufResponse) SetContent(data any) {
	protobufResponse.data = data
}

// Send sends the response
func (protobufResponse *ProtobufResponse) Send(w http.ResponseWriter, r *http.Request) {
	protobufBytes, err := proto.Marshal(protobufResponse.data.(proto.Message))
	if err != nil {
		panic(err)
	}
	protobufResponse.content = protobufBytes
	protobufResponse.SetHeader("Content-Type", mimetype.Protobuf)
	protobufResponse.Response.Send(w, r)
}
