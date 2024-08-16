package processing

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rmind-load-fbase/internal/model"
	"strconv"
)

func ProcessMessageData(client *firestore.Client, ctx context.Context, cloudEvent model.CloudEvent) error {
	data := model.MessageData{}
	dataBytes, err := json.Marshal(cloudEvent.Data)
	if err != nil {
		log.Printf("Failed to marshal cloudEvent.Data: %v", err)
		return err
	}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}

	collectionName := strconv.FormatInt(data.BoardID, 10)
	firestoreData := map[string]interface{}{
		"specversion": cloudEvent.SpecVersion,
		"id":          cloudEvent.ID,
		"source":      cloudEvent.Source,
		"type":        cloudEvent.Type,
		"time":        cloudEvent.Time,
		"data":        data,
		"object_key":  cloudEvent.ObjectKey,
	}

	_, err = client.Collection("messages").Doc("message").Collection(collectionName).Doc(strconv.FormatInt(data.ID, 10)).Set(ctx, firestoreData)
	if err != nil {
		log.Printf("Failed to save data to Firestore for BoardID %d: %v", data.BoardID, err)
		return err
	}

	fmt.Printf("Data successfully saved to Firestore for object: %s, BoardID: %d\n", cloudEvent.ObjectKey, data.BoardID)
	return nil
}
