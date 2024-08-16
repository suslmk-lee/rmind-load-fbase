package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"project-root/internal/processing"
	"project-root/pkg/s3"
	"time"
)

func ProcessObjectKey(sess *session.Session, client *firestore.Client, ctx context.Context, prefix string, key string) {
	for attempt := 1; attempt <= 3; attempt++ {
		cloudEvent, err := s3.ReadCloudEventFromS3(sess, bucketName, key)
		if err != nil {
			log.Printf("Attempt %d: Failed to read CloudEvent from S3: %v", attempt, err)
			if attempt == maxRetries {
				return
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		var processErr error
		switch prefix {
		case "rmine_push_data/messages":
			processErr = processing.ProcessMessageData(client, ctx, cloudEvent)
		case "rmine_push_data/users":
			processErr = processing.ProcessUserData(client, ctx, cloudEvent)
		case "rmine_push_data/issues":
			processErr = processing.ProcessIssueData(client, ctx, cloudEvent)
		default:
			log.Printf("Unknown prefix: %s", prefix)
			return
		}

		if processErr != nil {
			log.Printf("Attempt %d: Failed to save data to Firestore: %v", attempt, processErr)
			if attempt == maxRetries {
				return
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		newKey := fmt.Sprintf("processed/%s", key)
		err = s3.MoveObject(sess, bucketName, key, newKey)
		if err != nil {
			log.Printf("Failed to move object to %s: %v", newKey, err)
			return
		}

		fmt.Printf("Data successfully saved to Firestore for object: %s\n", key)
		break
	}
}
