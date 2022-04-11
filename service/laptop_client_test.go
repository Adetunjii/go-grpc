package service

import (
	"context"
	"github.com/Adetunjii/go-grpc/pb"
	"github.com/Adetunjii/go-grpc/sample"
	"github.com/Adetunjii/go-grpc/serializer"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id

	req := &pb.CreatelaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, expectedId, res.Id)

	other, err := laptopServer.Store.FindById(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	laptop1, err := serializer.ProtobufToJSON(laptop)
	require.NoError(t, err)

	laptop2, err := serializer.ProtobufToJSON(other)
	require.NoError(t, err)

	require.Equal(t, laptop1, laptop2)
}

func newTestLaptopClient(t *testing.T, address string) pb.LaptopServiceClient {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(connection)
}

func startTestLaptopServer(t *testing.T) (*LaptopServer, string) {
	laptopServer := NewLaptopServer(NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	//always run this in a seperate go routine, it's a blocking call
	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}
