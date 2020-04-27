package delivery

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/fileserver/proto/filetransfer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"io/ioutil"
	"log"
	"time"
)

type FileTransferDelivery struct {
}

func NewFileTransferDelivery() *FileTransferDelivery {
	return &FileTransferDelivery{}
}

func Interceptor(req interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()

	md, _ := metadata.FromIncomingContext(stream.Context())

	err := handler(req, stream)

	fmt.Printf(`--
	after incoming call=%v
	req=%#v
	time=%v
	md=%v
	err=%v`,
		info.FullMethod, req, time.Since(start), md, err)

	return err
}

func (uc *FileTransferDelivery) Upload(inStream filetransfer.UploadService_UploadServer) error {
	data := make([]byte, 0, 1024)
	md, _ := metadata.FromIncomingContext(inStream.Context())
	fileName := md.Get("fileName")[0]

	for {
		inData, err := inStream.Recv()
		if err == io.EOF {
			out := &filetransfer.UploadStatus{
				Message: "OK",
				Code:    filetransfer.UploadStatusCode_Ok,
			}
			fmt.Println("Transfer Ended")
			fmt.Printf("Filesize = %v", len(data))
			if err := inStream.SendAndClose(out); err != nil {
				log.Println(err)
			}
			break
		}
		if err != nil {
			return err
		}

		data = append(data, inData.Content...)
	}
	err := ioutil.WriteFile("resources/"+fileName, data, 0666)
	if err != nil {
		fmt.Println("fuck, error", err)
		return err
	}
	return nil
}

//func kek(fileName, fileType string) error {
//	filePath := filepath.Join(os.Getenv("FILE_ROOT")+ur.avatarDir, fileName+"."+fileType)
//
//	newFile, err := os.Create(filePath)
//	if err != nil {
//		return fmt.Errorf("failed to create file: %v", err)
//	}
//	defer newFile.Close()
//
//	_, err = io.Copy(newFile, file)
//	if err != nil {
//		return fmt.Errorf("error while writing to file: %v", err)
//	}
//}
