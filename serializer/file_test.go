package serializer

import (
	"github.com/Adetunjii/go-grpc/pb"
	"github.com/Adetunjii/go-grpc/sample"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWriteProtobufToBinaryFile(t *testing.T) {
	// allow tests to be run in parallel
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err  = ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	err = WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)
}

