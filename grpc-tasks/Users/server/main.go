package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "github.com/Ronak-Searce/grpc-tasks/users/proto"
	"github.com/Ronak-Searce/grpc-tasks/users/store"
	"google.golang.org/grpc"
)

const (
	port  = ":50051"
	dbURI = "projects/test-project-id/instances/test-instance/databases/crud1"
)

var users []*pb.UserInfo

type userServer struct {
	pb.UnimplementedUsererviceServer
}

func main() {
	store.SetUpSpanner(dbURI)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterUsererviceServer(s, &userServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func initUsers() {
	user1 := &pb.UserInfo{Id: "1", FirstName: "Ronak", LastName: "Babu"}
	user2 := &pb.UserInfo{Id: "2", FirstName: "Sanjeeb", LastName: "Kumar"}
	users = append(users, user1)
	users = append(users, user2)

}

func (s *userServer) GetUser(ctx context.Context, in *pb.Id) (*pb.UserInfo, error) {
	log.Printf("Received: %v", in)

	user, err := store.GetUser(in.Value, dbURI)
	if err != nil {
		return &pb.UserInfo{}, err
	}

	res := &pb.UserInfo{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	return res, nil
}

func (s *userServer) CreatUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.UserInfo, error) {
	log.Printf("Received: %v", in)

	newUser := store.UserInfo{
		FirstName: in.Firstname,
		LastName:  in.Lastname,
		Id:        strconv.Itoa(rand.Intn(100000000)),
	}

	user, err := store.CreateUser(dbURI, newUser)

	if err != nil {
		return &pb.UserInfo{}, err
	}
	res := &pb.UserInfo{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	return res, nil
}

func (s *userServer) UpdateUser(ctx context.Context,
	in *pb.UserInfo) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	user := store.UserInfo{
		Id:        in.Id,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}

	err := store.UpdateUser(dbURI, user)
	if err != nil {
		return &pb.Status{Value: int32(-1)}, err
	} else {
		return &pb.Status{Value: int32(1)}, nil
	}

}

func (s *userServer) DeleteUser(ctx context.Context,
	in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	err := store.DeleteUser(dbURI, in.Value)
	if err != nil {
		return &pb.Status{Value: int32(-1)}, err
	} else {
		return &pb.Status{Value: int32(1)}, nil
	}
}
