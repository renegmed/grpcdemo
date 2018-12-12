package main

import (
	"fmt"
	"grpc-demo/pb"
	"log"
	"net"

	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	log.Println("Starting server on port " + port)
	s.Serve(lis)
}

type employeeServer struct{}

// these methods are from message.pb.go type EmployeeServiceServer interface

func (s *employeeServer) GetByBadgeNumber(ctx context.Context,
	req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	fmt.Println("Request GetByBadgeNumber() is called.")

	if md, ok := metadata.FromContext(ctx); ok {
		fmt.Printf("Metadata received: %v\n", md)
	}

	return nil, nil
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
