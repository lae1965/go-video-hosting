package grpcclient

import (
	"bufio"
	"context"
	"go-video-hosting/gRPC/proto"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type FilesGRPCClient struct {
	gprcServer *FilesGRPCServer
}

var chunkSize int64

func NewFilesGRPCClient(gprcServer *FilesGRPCServer) *FilesGRPCClient {
	chunkSize = viper.GetInt64("gRPC.chunksize")

	return &FilesGRPCClient{
		gprcServer: gprcServer,
	}
}

func (client *FilesGRPCClient) SendToGRPCServer(ctx context.Context, fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		logrus.Errorf("Ошибка открытия файла: %s", err.Error())
		return "", err
	}
	defer file.Close()

	stream, err := client.gprcServer.Client.SendToGRPCServer(ctx)
	if err != nil {
		logrus.Errorf("Ошибка создания потока: %s", err.Error())
		return "", err
	}

	if err := stream.Send(&proto.FileSendRequest{
		Request: &proto.FileSendRequest_FileName{
			FileName: fileName,
		},
	}); err != nil {
		logrus.Errorf("Ошибка отправки имени файла: %s", err.Error())
		return "", err
	}

	reader := bufio.NewReader(file)
	buf := make([]byte, chunkSize)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			logrus.Errorf("Ошибка чтения файла: %s", err.Error())
			return "", err
		}

		if err := stream.Send(&proto.FileSendRequest{
			Request: &proto.FileSendRequest_Chunk{
				Chunk: buf[:n],
			},
		}); err != nil {
			logrus.Errorf("Ошибка отправки чанка: %s", err.Error())
			return "", err
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		logrus.Errorf("Ошибка закрытия потока: %s", err.Error())
		return "", err
	}

	return response.GetNewFileName(), nil
}

func (client *FilesGRPCClient) DeleteFromGRPCServer(ctx context.Context, fileName string) error {
	_, err := client.gprcServer.Client.DeleteFromGRPCServer(ctx, &proto.FileDeleteRequest{
		FileName: fileName,
	})

	return err
}

func (client *FilesGRPCClient) GetFromGRPCServer(ctx context.Context, fileName string, sendChank func(int64, string, []byte) error) error {
	var fileSize int64
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))

	stream, err := client.gprcServer.Client.GetFromGRPCServer(ctx, &proto.FileGetRequest{
		FileName: fileName,
		Start:    0,
		End:      -1,
	})
	if err != nil {
		logrus.Errorf("Ошибка создания потока: %s", err.Error())
		return err
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			logrus.Errorf("failed to receive a request: %s", err.Error())
			return err
		}

		if size := response.GetFileSize(); size > 0 {
			fileSize = size
		} else if fileSrteam := response.GetFileStream(); fileSrteam != nil {
			sendChank(fileSize, mimeType, fileSrteam)
		}

	}

	return nil
}
