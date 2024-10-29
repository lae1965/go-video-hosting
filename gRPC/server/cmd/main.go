package main

import (
	"go-video-hosting/gRPC/proto"
	"net"

	"go-video-hosting/gRPC/server"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	listen, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	srv := &grpcserver.FilesGRPCService{}
	proto.RegisterFilesServiceServer(server, srv)

	if err := server.Serve(listen); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
