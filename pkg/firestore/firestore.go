package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

func SaveData(client *firestore.Client, ctx context.Context, collection, doc string, data map[string]interface{}) error {
	_, err := client.Collection(collection).Doc(doc).Set(ctx, data)
	if err != nil {
		log.Printf("Failed to save data to Firestore: %v", err)
	}
	return err
}
