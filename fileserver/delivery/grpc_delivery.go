package delivery

import (
	"github.com/2020_1_no_homomorphism/fileserver/proto/filetransfer"
	"google.golang.org/grpc/metadata"
	"io"
	"io/ioutil"
	"log"
)

type FileTransferDelivery struct{}

func NewFileTransferDelivery() *FileTransferDelivery {
	return &FileTransferDelivery{}
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
			log.Println("Transfer Ended")
			log.Printf("Filesize = %v", len(data))
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
		log.Println("Error while saving file: ", err)
		return err
	}
	return nil
}
