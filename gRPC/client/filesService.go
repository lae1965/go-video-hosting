package grpcclient

import (
	"bufio"
	"context"
	"go-video-hosting/gRPC/proto"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type FilesGRPCClient struct {
	gprcServer *FilesGRPCServer
}

var chunkSize int64

func New(gprcServer *FilesGRPCServer) *FilesGRPCClient {
	chunkSize = viper.GetInt64("gRPC.chunksize")

	return &FilesGRPCClient{
		gprcServer: gprcServer,
	}
}

func (client *FilesGRPCClient) SendToGRPCServer(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
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
			FileName: fileHeader.Filename,
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

func (client *FilesGRPCClient) GetFromGRPCServer(ctx context.Context, fileName string, start, end int64, sendChunk func(int64, string, []byte) error) error {
	var fileSize int64
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))

	stream, err := client.gprcServer.Client.GetFromGRPCServer(ctx, &proto.FileGetRequest{
		FileName: fileName,
		Start:    start,
		End:      end,
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
			sendChunk(fileSize, mimeType, fileSrteam)
		}
	}

	return nil
}
