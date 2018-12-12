package main

import (
	"errors"
	"grpc-demo/pb"
	"log"
	"net"

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
	pb.RegisterEmployeeServiceServer(s, new(employeeServer))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Starting server on port %s\n" + port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type employeeServer struct{}

// these methods are from message.pb.go type EmployeeServiceServer interface

func (s *employeeServer) GetByBadgeNumber(ctx context.Context,
	req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	//fmt.Printf("Request GetByBadgeNumber() is called. Request badgenumber %d\n", req.BadgeNumber)

	for _, e := range employees {
		if req.BadgeNumber == e.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}

	return nil, errors.New("Employee not founde")
}

func (s *employeeServer) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	return nil
}

func (s *employeeServer) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeServer) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	return nil
}

func (s *employeeServer) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	return nil
}
