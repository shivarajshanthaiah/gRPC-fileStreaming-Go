package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	proto "github.com/shivarajshanthaiah-gRPC-sample-2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.StreamUploadClient

func main() {

	//connection to internal gRPC server
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewStreamUploadClient(conn)

	mb := 1024 * 1024 * 2
	uploadStreamFile("./11mb.txt", mb)
}

func uploadStreamFile(path string, batchSize int) {
	t := time.Now()
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, batchSize)
	batchNumber := 1
	stream, err := client.Upload(context.TODO())
	if err != nil {
		panic(err)
	}

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		chunk := buf[:num]
		if err := stream.Send(&proto.UploadRequest{FilePath: path, Chunks: chunk}); err != nil {
			fmt.Println(err)
			return
		}
		log.Printf("Sent - batch #%v - size - %v\n", batchNumber, len(chunk))
		batchNumber += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(t))
	log.Printf("Sent - %v bytes - %s\n", res.GetFileSize(), res.GetMessage())
}
