package s3

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rmind-load-fbase/common"
	"rmind-load-fbase/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	region    string
	endpoint  string
	accessKey string
	secretKey string
)

func init() {
	region = common.ConfInfo["nhn.region"]
	endpoint = common.ConfInfo["nhn.storage.endpoint.url"]
	accessKey = common.ConfInfo["nhn.storage.accessKey"]
	secretKey = common.ConfInfo["nhn.storage.secretKey"]
}

// CreateS3Session AWS S3 세션을 생성합니다.
func CreateS3Session() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true), // 경로 스타일을 강제 설정
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // TLS 검증 비활성화
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
		return nil, err
	}
	return sess, nil
}

// ListObjectsInS3Folder S3 폴더 내의 모든 객체를 목록화합니다.
func ListObjectsInS3Folder(sess *session.Session, bucketName, folder string) ([]string, error) {
	svc := s3.New(sess)

	var objectKeys []string
	err := svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folder),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, item := range page.Contents {
			objectKeys = append(objectKeys, *item.Key)
		}
		return !lastPage
	})

	if err != nil {
		return nil, err
	}
	return objectKeys, nil
}

// ReadObjectFromS3 S3 버킷에서 객체를 읽어옵니다.
func ReadObjectFromS3(sess *session.Session, bucketName, objectKey string) ([]byte, error) {
	// S3 서비스 클라이언트 생성
	svc := s3.New(sess)

	// S3 객체 읽기
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	// 객체 데이터 읽기
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ReadCloudEventFromS3 S3에서 CloudEvent 형식의 객체를 읽어옵니다.
func ReadCloudEventFromS3(sess *session.Session, bucketName, objectKey string) (model.CloudEvent, error) {
	body, err := ReadObjectFromS3(sess, bucketName, objectKey)
	if err != nil {
		return model.CloudEvent{}, err
	}

	var cloudEvent model.CloudEvent
	err = json.Unmarshal(body, &cloudEvent)
	if err != nil {
		return model.CloudEvent{}, err
	}

	// objectKey를 CloudEvent 구조체에 추가
	cloudEvent.ObjectKey = objectKey

	return cloudEvent, nil
}
