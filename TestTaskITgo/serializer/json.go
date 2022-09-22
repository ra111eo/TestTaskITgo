package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/runtime/protoiface"
)

func ProtobufToJSON(message protoiface.MessageV1) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "	",
		OrigName:     false,
	}
	return marshaler.MarshalToString(message)
}
