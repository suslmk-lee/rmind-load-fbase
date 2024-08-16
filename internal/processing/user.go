package processing

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"rmind-load-fbase/internal/model"
	"rmind-load-fbase/pkg/util"
	"strconv"
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

	err = util.ChangeTime(
		&data.CreatedOn,
		&data.UpdatedOn,
		&data.LastLoginOn.Time,
		&data.PasswdChangedOn.Time,
	)

	if err != nil {
		log.Printf("Failed to change data: %v", err)
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
	_, err = client.Collection("users").Doc(strconv.FormatInt(data.ID, 10)).Set(ctx, firestoreData)
	return err
}
