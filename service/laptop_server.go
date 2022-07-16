package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/Adetunjii/go-grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

type LaptopServer struct {
	LaptopStore LaptopStore
	ImageStore  ImageStore
	*pb.UnimplementedLaptopServiceServer
}

func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{
		LaptopStore: laptopStore,
		ImageStore:  imageStore,
	}
}

const maxImageSize = 1 << 20

// CreateLaptop
// Unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreatelaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("recieved a create laptop request wit id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		// generate new id
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop id: %v", err)
		}
		laptop.Id = id.String()
	}

	//time.Sleep(6 * time.Second)

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	err := server.LaptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, DuplicateException) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store %v", err)
	}
	log.Printf("saved laptop with id: %s", laptop.Id)
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return res, nil
}

func contextError(ctx context.Context) error {
	if ctx.Err() == context.Canceled {
		log.Print("cancelled")
		return status.Error(codes.Canceled, "cancelled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("error timeout")
		return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	return nil
}

func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {

	filter := req.GetFilter()
	log.Printf("received a search laptop filter with %v", filter)

	err := server.LaptopStore.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{Laptop: laptop}

		err := stream.Send(res)
		if err != nil {
			return err
		}

		log.Printf("sent laptop with id %s", laptop.GetId())
		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}

func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request for laptop %s with image type %s", laptopID, imageType)

	laptop, err := server.LaptopStore.FindById(laptopID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot find laptop: %v", err))
	}

	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "laptop %s doesn't exist", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {

		if err = contextError(stream.Context()); err != nil {
			return err
		}

		log.Print("waiting for data chunks")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}

		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	imageID, err := server.ImageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response %v", err))
	}

	log.Printf("image successfully saved %s, %d", imageID, imageSize)
	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}

	return err
}
