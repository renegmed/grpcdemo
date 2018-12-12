package main

import (
	"flag"
	"fmt"
	"grpc-demo/pb"
	"log"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const port = ":9000"

func main() {
	log.Println("Start....")
	option := flag.Int("o", 1, "Command to run")
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Got credential....")
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)

	//log.Println("Option %v", *option)
	switch *option {
	case 1:
		SendMetadata(client)
	default:
		log.Println("Option is not valid.")
	}
}

func SendMetadata(client pb.EmployeeServiceClient) {

	log.Println("Start SendMetadata....")
	md := metadata.MD{}
	md["user"] = []string{"mvansickle"}
	md["password"] = []string{"password1"}
	ctx := context.Background()
	resp, err := client.GetByBadgeNumber(ctx, &pb.GetByBadgeNumberRequest{BadgeNumber: 5})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Response received: %v\n", resp)
}
