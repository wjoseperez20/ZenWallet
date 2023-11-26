package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
	"time"
)

var Aws *session.Session

func ConnectAWS() {
	var awsSession *session.Session
	var err error

	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	for i := 1; i <= 3; i++ {
		awsSession, err = session.NewSession(&aws.Config{
			Region:      aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
			Endpoint:    aws.String("s3.amazonaws.com"), // Specify the S3 endpoint
		})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize AWS session. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	Aws = awsSession
}
