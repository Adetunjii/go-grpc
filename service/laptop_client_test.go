package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Adetunjii/go-grpc/pb"
	"github.com/Adetunjii/go-grpc/sample"
	"github.com/Adetunjii/go-grpc/serializer"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := NewInMemoryLaptopStore()
	laptopServer, serverAddress := startTestLaptopServer(t, laptopStore, nil)
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

	other, err := laptopServer.LaptopStore.FindById(res.Id)
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

func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	testimageFolder := "../tmp"

	laptopStore := NewInMemoryLaptopStore()
	imageStore := NewDiskImageStore(testimageFolder)

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	_, serverAddress := startTestLaptopServer(t, laptopStore, imageStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/k-mean-algorithm.jpg", testimageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptop.GetId(),
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}
	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, size, res.GetSize())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s/%s%s", testimageFolder, res.GetId(), imageType)
	require.FileExists(t, savedImagePath)
	require.NoError(t, os.Remove(savedImagePath))
}

func startTestLaptopServer(t *testing.T, laptopStore LaptopStore, imageStore ImageStore) (*LaptopServer, string) {
	laptopServer := NewLaptopServer(laptopStore, imageStore)

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	//always run this in a seperate go routine, it's a blocking call
	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCupGhb:   2.2,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	laptopStore := NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberOfCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberOfCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberOfCores = 4
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = 5.5
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}

		err := laptopStore.Save(laptop)
		require.NoError(t, err)
	}

	_, serverAddress := startTestLaptopServer(t, laptopStore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())

		found += 1
	}

	require.Equal(t, len(expectedIDs), found)
}
