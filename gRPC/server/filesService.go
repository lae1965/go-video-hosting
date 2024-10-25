package grpcserver

import (
	"context"
	"fmt"
	"go-video-hosting/gRPC/proto"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FilesGRPCService struct {
	proto.UnimplementedFilesServiceServer
}

func getFullFileName(fileName string) string {
	return filepath.Join("gRPC", "server", "data", fileName)
}

func generateUniqueFileName(fileName string) string {
	var uniqueFileName string
	for {
		uniqueFileName = fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fileName))
		if _, err := os.Stat(getFullFileName(uniqueFileName)); os.IsNotExist(err) {
			break
		}
	}
	return uniqueFileName
}

func (s *FilesGRPCService) SendToGRPCServer(stream proto.FilesService_SendToGRPCServerServer) error {
	var file *os.File
	var responseError error
	var wg sync.WaitGroup
	var fileName string

	chunkChan := make(chan []byte)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for chunk := range chunkChan {
			if file != nil {
				if _, err := file.Write(chunk); err != nil {
					responseError = status.Errorf(codes.Internal, "failed to write chunk to file: %s", err.Error())
					return
				}
			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			responseError = status.Errorf(codes.Internal, "failed to receive a request: %s", err.Error())
			break
		}

		switch request := req.Request.(type) {
		case *proto.FileSendRequest_FileName:
			fileName = generateUniqueFileName(request.FileName)
			file, err = os.Create(getFullFileName(fileName))
			if err != nil {
				responseError = status.Errorf(codes.Internal, "failed to create file: %s", err.Error())
			}

		case *proto.FileSendRequest_Chunk:
			if file == nil {
				responseError = status.Errorf(codes.InvalidArgument, "file has not been initialized")
			}

			chunkChan <- request.Chunk

		default:
			responseError = status.Errorf(codes.Internal, "unknown request type: %T", request)
		}

		if responseError != nil {
			break
		}
	}

	close(chunkChan)
	wg.Wait()

	if file != nil {
		file.Close()
	}

	if err := stream.SendAndClose(&proto.FileSendResponse{
		NewFileName: fileName,
	}); err != nil {
		responseError = status.Errorf(codes.Internal, "failed to close stream: %s", err.Error())
	}

	if responseError != nil {
		logrus.Printf("%s", responseError.Error())
	}

	return responseError
}

func (s *FilesGRPCService) DeleteFromGRPCServer(ctx context.Context, request *proto.FileDeleteRequest) error {
	return os.Remove(getFullFileName(request.GetFileName()))
}

func (s *FilesGRPCService) GetFromGRPCServer(req *proto.FileGetRequest, stream proto.FilesService_GetFromGRPCServerServer) error {
	file, err := os.Open(getFullFileName(req.GetFileName()))
	if err != nil {
		return status.Errorf(codes.Internal, "failed to open file: %s", err.Error())
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to stat file: %s", err.Error())
	}

	fileSize := fileInfo.Size()

	// Отправляем размер файла
	if err := stream.Send(&proto.FileGetResponse{
		Response: &proto.FileGetResponse_FileSize{
			FileSize: fileSize,
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to send file size: %s", err.Error())
	}

	// Определяем отправлять весь файл или фрагмент
	start := req.GetStart()
	end := req.GetEnd()

	if end == -1 {
		start = 0
		end = fileSize
	} else if start < 0 || end > fileSize || start > end {
		return status.Errorf(codes.Internal, "invalid start or end index: start=%d, end=%d", start, end)
	}

	if _, err := file.Seek(start, 0); err != nil {
		return status.Errorf(codes.Internal, "failed to seek in file: %s", err.Error())
	}

	//Зачитываем и отправляем чанками
	chunkSize := int64(1024)
	for pos := start; pos < end; {
		sizeToRead := chunkSize
		if pos+sizeToRead > end {
			sizeToRead = end - pos
		}

		chunk := make([]byte, sizeToRead)
		n, err := file.Read(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to read file chunk: %s", err.Error())
		}

		if n > 0 {
			if err := stream.Send(&proto.FileGetResponse{
				Response: &proto.FileGetResponse_FileStream{
					FileStream: chunk[:n],
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to send file chunk: %s", err.Error())
			}

			pos += int64(n)
		}
	}

	return nil
}
