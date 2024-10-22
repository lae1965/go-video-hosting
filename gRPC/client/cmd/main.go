package main

import (
	// "bufio"
	"context"
	// "fmt"
	"go-video-hosting/gRPC/proto"
	"log"
	// "os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const chunkSize = 1024

func main() {
	// Устанавливаем соединение с сервером
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	// Создаем клиента
	client := proto.NewFilesServiceClient(conn)

	// // Открываем файл
	// filePath := "photo_2023-11-19_22-27-01.jpg" // замените на ваш файл
	// file, err := os.Open(filePath)
	// if err != nil {
	// 	log.Fatalf("Ошибка открытия файла: %v", err)
	// }
	// defer file.Close()

	// // Начинаем потоковую передачу
	// stream, err := client.SendToGRPCServer(context.Background())
	// if err != nil {
	// 	log.Fatalf("Ошибка создания потока: %v", err)
	// }

	// // Отправляем имя файла
	// err = stream.Send(&proto.FileSendRequest{Request: &proto.FileSendRequest_FileName{FileName: file.Name()}})
	// if err != nil {
	// 	log.Fatalf("Ошибка отправки имени файла: %v", err)
	// }

	// // Читаем файл по частям и отправляем чанками
	// reader := bufio.NewReader(file)
	// buffer := make([]byte, chunkSize) // Размер чанка - 1024 байта

	// for {
	// 	n, err := reader.Read(buffer)
	// 	if err != nil {
	// 		if err.Error() != "EOF" {
	// 			log.Fatalf("Ошибка чтения файла: %v", err)
	// 		}
	// 		log.Println("EOF")
	// 		break
	// 	} else {
	// 		log.Println("Reading...")
	// 	}

	// 	// Отправляем чанки
	// 	err = stream.Send(&proto.FileSendRequest{
	// 		Request: &proto.FileSendRequest_Chunk{
	// 			Chunk: buffer[:n],
	// 		},
	// 	})
	// 	if err != nil {
	// 		log.Fatalf("Ошибка отправки чанка: %v", err)
	// 	}
	// }

	// // Завершаем потоковую передачу и получаем ответ
	// response, err := stream.CloseAndRecv()
	// if err != nil {
	// 	log.Fatalf("Ошибка получения ответа: %v", err)
	// }

	// log.Printf("response: %v", response)

	// fmt.Printf("Статус ответа: %v\n", response.Status)
	// fmt.Printf("Имя файла: %v\n", response.NewFileName)

	response, err := client.DeleteFromGRPCServer(context.Background(), &proto.FileStandartRequest{FileName: "5af5f5b1-f870-4d27-876c-bafbf4dc366f.jpg"})

	log.Printf("response: %v", response.GetStatus())
	log.Printf("err: %v", err)

}
