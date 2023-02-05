package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Ronak-Searce/grpc-tasks/users/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUsererviceClient(conn)

	// runGetUser(client, "98498081")
	runCreateUser(client, "21", "Ronak", "Babu")
	// runUpdateUser(client, "98498081", "Ronak", "Babu")
	// runDeleteUser(client, "98498081")
}

func runGetUser(client pb.UsererviceClient, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: userid}
	res, err := client.GetUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetUser(_) = _, %v", client, err)
	}
	log.Printf("UserInfo: %v", res)
}

func runCreateUser(client pb.UsererviceClient, userid string, first_name string,
	last_name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.UserInfo{FirstName: first_name, LastName: last_name, Id: userid}
	res, err := client.CreatUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateUser(_) = _, %v", client, err)
	}
	if res.String() != "" {
		log.Printf("User created: %v", res)
	} else {
		log.Printf("CreateUser Failed")
	}

}

func runUpdateUser(client pb.UsererviceClient, userid string,
	first_name string, last_name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.UserInfo{Id: userid, FirstName: first_name, LastName: last_name}
	res, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdateUser(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("UpdateUser Success")
	} else {
		log.Printf("UpdateUser Failed")
	}
}

func runDeleteUser(client pb.UsererviceClient, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: userid}
	res, err := client.DeleteUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteUser(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("DeleteUser Success")
	} else {
		log.Printf("DeleteUser Failed")
	}
}
