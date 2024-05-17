package response

import (
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

var fileTestProtoMsgTypes = make([]protoimpl.MessageInfo, 1)

type data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key1 string `protobuf:"bytes,1,opt,name=key1,proto3" json:"key1,omitempty"`
	Key2 string `protobuf:"bytes,2,opt,name=key2,proto3" json:"key2,omitempty"`
}

func (x *data) Reset() {
	*x = data{}
	if protoimpl.UnsafeEnabled {
		mi := &fileTestProtoMsgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*data) ProtoMessage() {}

func (x *data) ProtoReflect() protoreflect.Message {
	mi := &fileTestProtoMsgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func TestProtobufResponse(t *testing.T) {
	protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(struct{}{}).PkgPath(),
			RawDescriptor: []byte{
				0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x65,
				0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12,
				0x0a, 0x04, 0x6b, 0x65, 0x79, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65,
				0x79, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
				0x52, 0x04, 0x6b, 0x65, 0x79, 0x32, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x72, 0x65, 0x73, 0x70,
				0x6f, 0x6e, 0x73, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
			},
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes: []interface{}{
			(*data)(nil), // 0: response.data
		},
		DependencyIndexes: []int32{
			0, // [0:0] is the sub-list for method output_type
			0, // [0:0] is the sub-list for method input_type
			0, // [0:0] is the sub-list for extension type_name
			0, // [0:0] is the sub-list for extension extendee
			0, // [0:0] is the sub-list for field type_name
		},
		MessageInfos: fileTestProtoMsgTypes,
	}.Build()

	d := new(data)
	d.Key1 = "value1"
	d.Key2 = "value2"
	response := NewResponse(200).Protobuf(d)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	response.Send(recorder, request)
	result := recorder.Result()
	body := result.Body
	content, err := io.ReadAll(body)
	assert.Nil(t, err)
	defer func() { body.Close() }()
	d2 := new(data)
	err = proto.Unmarshal(content, d2)
	assert.Nil(t, err)
	assert.Equal(t, "value1", d2.Key1)
	assert.Equal(t, "value2", d2.Key2)
}
