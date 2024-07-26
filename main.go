package main

import (
	"context"
	"fmt"
	"log"
	"rmind-load-fbase/common"
	"rmind-load-fbase/s3"
	"sync"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	bucketName     string
	firestoreCreds string
	objectPrefix   string // S3 버킷의 객체 접두사
)

func init() {
	bucketName = common.ConfInfo["nhn.storage.bucket.name"]
	firestoreCreds = common.ConfInfo["firestore.cred.file"]
	objectPrefix = common.ConfInfo["firestore.object.prefix"] // S3 객체의 경로
}

func main() {
	// AWS S3 세션 설정
	sess, err := s3.CreateS3Session()
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// S3 버킷에서 객체 목록 가져오기
	objectKeys, err := s3.ListObjectsInBucket(sess, bucketName, objectPrefix)
	if err != nil {
		log.Fatalf("Failed to list objects in S3 bucket: %v", err)
	}

	// Firebase Firestore 초기화
	ctx := context.Background()
	sa := option.WithCredentialsFile(firestoreCreds)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// 병렬 처리를 위한 WaitGroup
	var wg sync.WaitGroup

	// 각 객체를 읽어와 Firestore에 저장
	for _, objectKey := range objectKeys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()

			// S3에서 CloudEvent 읽기
			cloudEvent, err := s3.ReadCloudEventFromS3(sess, bucketName, key)
			if err != nil {
				log.Printf("Failed to read CloudEvent from S3: %v", err)
				return
			}

			// Firestore에 데이터 저장
			data := map[string]interface{}{
				"specversion":     cloudEvent.SpecVersion,
				"id":              cloudEvent.ID,
				"source":          cloudEvent.Source,
				"type":            cloudEvent.Type,
				"datacontenttype": cloudEvent.DataContentType,
				"time":            cloudEvent.Time,
				"data":            cloudEvent.Data,
				"object_key":      cloudEvent.ObjectKey,
			}

			_, err = client.Collection("your_collection").Doc(cloudEvent.ID).Set(ctx, data)
			if err != nil {
				log.Printf("Failed to save data to Firestore: %v", err)
			} else {
				fmt.Printf("Data successfully saved to Firestore for object: %s\n", key)
			}
		}(objectKey)
	}

	// 모든 고루틴이 완료될 때까지 대기
	wg.Wait()
}
