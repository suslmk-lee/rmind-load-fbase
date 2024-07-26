package main

import (
	"context"
	"fmt"
	"log"
	"rmind-load-fbase/common"
	"rmind-load-fbase/s3"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	bucketName        string
	firestoreCreds    string
	objectKey         string
	lastProcessedTime time.Time
)

func init() {
	bucketName = common.ConfInfo["nhn.storage.bucket.name"]
	firestoreCreds = common.ConfInfo["firestore.cred.file"]
	objectKey = common.ConfInfo["firestore.object.key"]
}

func main() {
	// AWS S3 세션 설정
	sess, err := s3.CreateS3Session()
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// S3에서 CloudEvent 형식의 데이터 읽기
	cloudEvent, err := s3.ReadCloudEventFromS3(sess, bucketName, objectKey)
	if err != nil {
		log.Fatalf("Failed to read CloudEvent from S3: %v", err)
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
		log.Fatalf("Failed to save data to Firestore: %v", err)
	}

	fmt.Println("Data successfully saved to Firestore")
}
