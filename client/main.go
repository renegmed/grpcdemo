package main

import (
	"flag"
	"fmt"
	"grpc-demo/pb"
	"io"
	"log"
	"os"

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
	case 2:
		GetByBadgeNumber(client)
	case 3:
		GetAll(client)
	case 4:
		AddPhoto(client)
	case 5:
		SaveAll(client)
	default:
		log.Println("Option is not valid.")
	}
}

func SaveAll(client pb.EmployeeServiceClient) {
	employees := []pb.Employee{
		pb.Employee{
			BadgeNumber:         123,
			FirstName:           "Joh",
			LastName:            "Smith",
			VacationAccrualRate: 1.2,
			VacationAccrued:     0,
		},
		pb.Employee{
			BadgeNumber:         234,
			FirstName:           "Lisa",
			LastName:            "Wu",
			VacationAccrualRate: 1.7,
			VacationAccrued:     10,
		},
	}
	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})

	// receiving employee data
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneCh <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res.Employee)
		}
	}()

	// Sending employee data
	for _, e := range employees {
		err := stream.Send(&pb.EmployeeRequest{Employee: &e})
		if err != nil {
			log.Fatal(err)
		}
	}
	stream.CloseSend()

	// wait until of all tasks are done e.g. sending and receiving data
	<-doneCh

}
func AddPhoto(client pb.EmployeeServiceClient) {
	f, err := os.Open("Penguins.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	md := metadata.New(map[string]string{"badgenumber": "2080"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for {
		chunk := make([]byte, 64*1024)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		stream.Send(&pb.AddPhotoRequest{Data: chunk})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.IsOk)
}

func GetAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.Employee)
	}
}

func GetByBadgeNumber(client pb.EmployeeServiceClient) {
	res, err := client.GetByBadgeNumber(context.Background(),
		&pb.GetByBadgeNumberRequest{BadgeNumber: 2080})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Employee)
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
