package main

import (
	"log"
	"os"

	"github.com/abdullah2993/goias3/ias3"
)

var (
	accessKey = os.Getenv("IAS3_ACCESS_KEY")
	secretKey = os.Getenv("IAS3_SECRET_KEY")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if accessKey == "" || secretKey == "" {
		log.Fatalf("unable to find the access/secret key")
	}

	f, err := os.Open("./uploads/file.mp4")
	if err != nil {
		log.Fatalf("unable to read file")
	}
	err = ias3.NewReq().
		AutoCreateBucket().
		WithMeta("title", "Test Upload").
		WithCreds(accessKey, secretKey).
		UploadFile("test_loaylty_exe_1/file.mp4", f)
	if err != nil {
		log.Fatalf("upload filed: %v", err)
	}

}
