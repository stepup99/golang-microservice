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
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	region := os.Getenv("KINESIS_REGION")
	if region == "" {
		log.Fatalf("KINESIS_REGION environment variable is not set")
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

	// Define the stream name and partition key to filter
	streamName := os.Getenv("KINESIS_STREAM_NAME")
	partitionKeyToFilter := "partitionKey-1"
	shardId := "shardId-000000000001" // Specific shard ID

	// Get the shard iterator
	shardIteratorResp, err := svc.GetShardIterator(&kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamName),
		ShardId:           aws.String(shardId),
		ShardIteratorType: aws.String("TRIM_HORIZON"), // Use "TRIM_HORIZON" to start from the oldest record
	})
	if err != nil {
		log.Fatalf("Failed to get shard iterator: %v", err)
	}

	shardIterator := shardIteratorResp.ShardIterator

	// Continuously read records from the stream
	for {
		getRecordsResp, err := svc.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
			Limit:         aws.Int64(10), // Adjust as needed
		})
		if err != nil {
			log.Fatalf("Failed to get records: %v", err)
		}

		// Filter and print records based on the partition key
		for _, record := range getRecordsResp.Records {
			if record.PartitionKey != nil && *record.PartitionKey == partitionKeyToFilter {
				fmt.Printf("Record: %s\n", string(record.Data))
			}
		}

		// Update the shard iterator and sleep before polling again
		shardIterator = getRecordsResp.NextShardIterator
		time.Sleep(1 * time.Second) // Adjust polling interval as needed
	}
}
