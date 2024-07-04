package gapi

import (
	context "context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/morgan/simplebank/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"time"
)

func GAPIClient() {
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimplebankClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// response, err := client.CreateUser(ctx, &pb.CreateUserRequest{
	// 	Username: "morgan123",
	// 	FullName: "Michael ade",
	// 	Email:    "michael@gmail.com",
	// 	Password: "password",
	// })
	response, err := client.LoginUser(ctx, &pb.LoginUserRequest{
		Username: "morgan123",
		Password: "password3",
	})
	if err != nil {
		log.Fatalf("CreateUser: %v", err)
	}
	byt, err := json.Marshal(response.GetUser())
	if err != nil {
		log.Fatalf("Error converting to json: %v", err)
	}
	fmt.Println(string(byt))
}
