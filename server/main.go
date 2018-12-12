package main

import (
	"errors"
	"fmt"
	"grpc-demo/pb"
	"io"
	"log"
	"net"

	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(opts...)
	pb.RegisterEmployeeServiceServer(s, new(employeeService))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Starting server on port %s\n" + port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type employeeService struct{}

// these methods are from message.pb.go type EmployeeServiceServer interface

func (s *employeeService) GetByBadgeNumber(ctx context.Context,
	req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	for _, e := range employees {
		if req.BadgeNumber == e.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}
	return nil, errors.New("Employee not founde")
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		stream.Send(&pb.EmployeeResponse{Employee: &e})
	}
	return nil
}

func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromContext(stream.Context())
	if ok {
		fmt.Printf("Receiving photo for badge number %v\n", md["badgenumber"][0])
	}
	imgData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File received with length: %v\n", len(imgData))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}

		fmt.Printf("Received %v bytes\n", (data.Data))

		imgData = append(imgData, data.Data...)
	}
	return nil
}
