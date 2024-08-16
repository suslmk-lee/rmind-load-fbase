package processing

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"rmind-load-fbase/internal/model"
)

func ProcessUserData(client *firestore.Client, ctx context.Context, cloudEvent model.CloudEvent) error {
	data := model.UserData{}
	dataBytes, err := json.Marshal(cloudEvent.Data)
	if err != nil {
		log.Printf("Failed to marshal cloudEvent.Data: %v", err)
		return err
	}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}
	firestoreData := map[string]interface{}{
		"specversion": cloudEvent.SpecVersion,
		"id":          cloudEvent.ID,
		"source":      cloudEvent.Source,
		"type":        cloudEvent.Type,
		"time":        cloudEvent.Time,
		"data":        data,
		"object_key":  cloudEvent.ObjectKey,
	}
	_, err = client.Collection("users").Doc(cloudEvent.ID).Set(ctx, firestoreData)
	return err
}
