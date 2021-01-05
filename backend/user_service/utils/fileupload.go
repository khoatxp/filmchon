package utils

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type fileUpload struct{}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader, userID string) (string, map[string]string)
}

//So what is exposed is Uploader
var FileUpload UploadFileInterface = &fileUpload{}

func (fu *fileUpload) UploadFile(file *multipart.FileHeader, userID string) (string, map[string]string) {

	errList := map[string]string{}

	f, err := file.Open()
	if err != nil {
		errList["Not_Image"] = "Please Upload a valid image"
		return "", errList
	}
	defer f.Close()

	size := file.Size
	//The image should not be more than 500KB
	fmt.Println("the size: ", size)
	if size > int64(512000) {
		errList["Too_large"] = "Sorry, Please upload an Image of 500KB or less"
		return "", errList

	}
	//only the first 512 bytes are used to sniff the content type of a file,
	//so no need to read the entire bytes of a file.
	buffer := make([]byte, size)
	f.Read(buffer)
	fileType := http.DetectContentType(buffer)
	//if the image is valid
	if !strings.HasPrefix(fileType, "image") {
		errList["Not_Image"] = "Please Upload a valid image"
		return "", errList
	}

	//Format file's path
	filePath := userID + "/" + "avatar" + path.Ext(file.Filename)

	//Get env values
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	//new session for uploading files
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secKey,
				"",
			),
		})

	if err != nil {
		log.Fatal(err)
	}

	//upload to the s3 bucket
	uploader := s3manager.NewUploader(sess)
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		ACL:         aws.String("public-read"),
		Key:         aws.String(filePath),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(http.DetectContentType(buffer)),
	})
	_ = up

	if err != nil {
		fmt.Println("error", err)
		errList["Other_Err"] = "something went wrong"
		return "", errList
	}

	fileURL := "https://" + bucketName + "." + "s3-" + region + ".amazonaws.com/" + filePath
	return fileURL, nil
}
