package grpcclient

import (
	"fmt"
	"go-video-hosting/gRPC/proto"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FilesGRPCServer struct {
	Connection *grpc.ClientConn
	Client     proto.FilesServiceClient
}

func NewFilesGRPCServer() (*FilesGRPCServer, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", viper.GetString("gRPC.host"), viper.GetString("gRPC.port")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &FilesGRPCServer{
		Connection: conn,
		Client:     proto.NewFilesServiceClient(conn),
	}, nil
}
