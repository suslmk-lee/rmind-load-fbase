package main

import (
	"context"
	"fmt"
	"log"
	"rmind-load-fbase/common"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	firestoreCreds string
)

func init() {
	firestoreCreds = common.ConfInfo["firestore.cred.file"]
}

type User struct {
	Specversion     string    `json:"specversion"`
	ID              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	Datacontenttype string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            Data      `json:"data"`
}

type Data struct {
	ID               int64        `json:"id"`
	Login            string       `json:"login"`
	HashedPassword   string       `json:"hashed_password"`
	Firstname        string       `json:"firstname"`
	Lastname         string       `json:"lastname"`
	Admin            bool         `json:"admin"`
	Status           int64        `json:"status"`
	LastLoginOn      On           `json:"last_login_on"`
	Language         string       `json:"language"`
	AuthSourceID     AuthSourceID `json:"auth_source_id"`
	CreatedOn        time.Time    `json:"created_on"`
	UpdatedOn        time.Time    `json:"updated_on"`
	Type             string       `json:"type"`
	MailNotification string       `json:"mail_notification"`
	Salt             string       `json:"salt"`
	MustChangePasswd bool         `json:"must_change_passwd"`
	PasswdChangedOn  On           `json:"passwd_changed_on"`
}

type AuthSourceID struct {
	Int64 int64 `json:"Int64"`
	Valid bool  `json:"Valid"`
}

type On struct {
	Time  time.Time `json:"Time"`
	Valid bool      `json:"Valid"`
}

func main() {
	// Firebase Admin SDK 초기화
	sa := option.WithCredentialsFile(firestoreCreds)
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Firestore 클라이언트 생성
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v\n", err)
	}
	defer client.Close()

	// Firestore에서 특정 컬렉션의 모든 문서를 가져오기
	getAllUsers(client)
}

// 특정 컬렉션의 모든 문서를 가져오는 함수
func getAllUsers(client *firestore.Client) {
	ctx := context.Background()

	// "users" 컬렉션의 모든 문서를 가져오기
	iter := client.Collection("users").Documents(ctx)
	defer iter.Stop()

	// 각 문서를 반복하며 출력
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		// 문서 데이터를 User 구조체로 변환
		var user User
		err = doc.DataTo(&user)
		if err != nil {
			log.Fatalf("Failed to convert document data: %v", err)
		}

		// 출력
		fmt.Printf("User ID: %s, Name: %s %s, Admin: %v\n", user.ID, user.Data.Firstname, user.Data.Lastname, user.Data.CreatedOn)
	}
}
