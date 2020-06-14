package main

import (
	"github.com/UstinovV/wm_api/auth"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("Error to start server ", err)
	}

	s := auth.AuthServer{}

	grpcServer := grpc.NewServer()

	auth.RegisterAuthCheckerServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve grpc ", err)
	}
}
