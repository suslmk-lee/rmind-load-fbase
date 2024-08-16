package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"rmind-load-fbase/common"
	"rmind-load-fbase/internal/processing"
	"rmind-load-fbase/pkg/s3"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	bucketName     string
	firestoreCreds string
	objectPrefixs  = []string{
		"rmine_push_data/messages",
		"rmine_push_data/users",
		"rmine_push_data/issues",
	} // S3 버킷의 객체 접두사
	maxRetries    = 3
	maxGoroutines = 5 // S3 버킷에서 객체 수 개수 제한 (3000개 이상은 제한 없음)
)

func init() {
	bucketName = common.ConfInfo["nhn.storage.bucket.name"]
	firestoreCreds = common.ConfInfo["firestore.cred.file"]
	//objectPrefix = common.ConfInfo["firestore.object.prefix"] // S3 객체의 경로
}

func processObjectKey(sess *session.Session, client *firestore.Client, ctx context.Context, prefix string, key string) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		cloudEvent, err := s3.ReadCloudEventFromS3(sess, bucketName, key)
		if err != nil {
			log.Printf("Attempt %d: Failed to read CloudEvent from S3: %v", attempt, err)
			if attempt == maxRetries {
				return
			}
			time.Sleep(time.Duration(attempt) * time.Second) // 재시도 전에 대기
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
			time.Sleep(time.Duration(attempt) * time.Second) // 재시도 전에 대기
			continue
		}

		// 처리 완료된 객체를 /processed 경로로 이동
		newKey := fmt.Sprintf("processed/%s", key)
		err = s3.MoveObject(sess, bucketName, key, newKey)
		if err != nil {
			log.Printf("Failed to move object to %s: %v", newKey, err)
			return
		}

		fmt.Printf("Data successfully saved to Firestore for object: %s\n", key)
		break // 성공하면 루프 탈출
	}
}

func main() {
	// AWS S3 세션 설정
	sess, err := s3.CreateS3Session()
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// Firebase Firestore 초기화
	ctx := context.Background()
	fmt.Println("Initializing Firestore...")
	sa := option.WithCredentialsFile(firestoreCreds)
	fmt.Println("Creating Firebase app...")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	sem := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	// 각 객체를 읽어와 Firestore에 저장
	for _, prefix := range objectPrefixs {
		fmt.Printf("Processing %s...\n", prefix)
		objectKeys, err := s3.ListObjectsInBucket(sess, bucketName, prefix)
		if err != nil {
			log.Printf("Failed to list objects for prefix %s: %v", prefix, err)
			continue
		}

		for _, objectKey := range objectKeys {
			wg.Add(1)
			fmt.Println("Processing object: ", objectKey)
			sem <- struct{}{}

			go func(pfx, key string) {
				defer wg.Done()
				defer func() { <-sem }()

				processObjectKey(sess, client, ctx, pfx, key)
			}(prefix, objectKey)
		}
	}

	// 모든 고루틴이 완료될 때까지 대기
	wg.Wait()

	fmt.Println("All tasks completed. Exiting...")
}
