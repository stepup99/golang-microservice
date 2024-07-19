package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Create a new session with hardcoded credentials and default region
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("KINESIS_REGION")), // Replace with your region
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create a new Kinesis client
	svc := kinesis.New(sess)

	// Define the stream name
	streamName := os.Getenv("KINESIS_STREAM_NAME")

	// Produce records to the Kinesis stream
	for i := 0; i < 10; i++ {
		data := fmt.Sprintf("data 11:26 AM %d", i)
		_, err := svc.PutRecord(&kinesis.PutRecordInput{
			Data:         []byte(data),
			PartitionKey: aws.String("partitionKey-1"),
			StreamName:   aws.String(streamName),
		})
		if err != nil {
			log.Fatalf("Failed to put record: %v", err)
		}
		fmt.Printf("Record %d successfully put into the stream\n", i)
		time.Sleep(1 * time.Second) // Sleep for a second to simulate time between records
	}
}
