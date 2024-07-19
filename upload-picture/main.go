package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var region string
var accessKey string
var secretKey string
var bucketName string
var uploader *s3manager.Uploader

func main() {
	fmt.Print("LLLLLLLLLLLLLL")

	r := gin.Default()
	r.POST("/upload", uploadFile)
	r.Run(":8080")
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName = os.Getenv("AWS_BUCKET")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("KINESIS_REGION")), // Replace with your region
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	uploader = s3manager.NewUploader(sess)

}

func uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var errors []string
	var uploadURLs []string
	files := form.File["files"]

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error opening file %s: %s ", file.Filename, err.Error()))
			continue
		}
		defer f.Close()
		uploadURL, err := saveFile(f, file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error saving file %s : %s ", file.Filename, err.Error()))
		}
		uploadURLs = append(uploadURLs, uploadURL)
	}
	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors})
	} else {
		c.JSON(http.StatusOK, gin.H{"url": uploadURLs})
	}
}

func saveFile(fileReader io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileHeader.Filename),
		Body:   fileReader,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileHeader.Filename)

	return url, nil
}
