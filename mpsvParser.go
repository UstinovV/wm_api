package main

import (
	"github.com/UstinovV/wm_api/mpsv"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":80011")
	if err != nil {
		log.Fatal("Error to start server ", err)
	}

	s := mpsv.MpsvServer{}

	grpcServer := grpc.NewServer()
	mpsv.RegisterMpsvParserServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve grpc " , err)
	}
}
