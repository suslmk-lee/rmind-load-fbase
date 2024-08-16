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

func ProcessIssueData(client *firestore.Client, ctx context.Context, cloudEvent model.CloudEvent) error {
	data := model.IssueData{}
	dataBytes, err := json.Marshal(cloudEvent.Data)
	if err != nil {
		log.Printf("Failed to marshal cloudEvent.Data: %v", err)
		return err
	} else {
		log.Printf("Processing IssueData:: %s", cloudEvent.ObjectKey)
	}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}

	data.Description = util.ChangeString(data.Description)
	data.Subject = util.ChangeString(data.Subject)
	data.Notes = util.ChangeString(data.Notes)

	err = util.ChangeTime(
		&data.DueDate,
		&data.CreatedOn,
		&data.UpdatedOn,
		&data.StartDate,
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
	_, err = client.Collection("issues").Doc(strconv.Itoa(data.ID)).Set(ctx, firestoreData)

	return err
}

func ProcessRawIssueData(client *firestore.Client, ctx context.Context, cloudEvent model.CloudEvent) error {
	data := model.RawIssueData{}
	dataBytes, err := json.Marshal(cloudEvent.Data)
	if err != nil {
		log.Printf("Failed to marshal cloudEvent.Data: %v", err)
		return err
	} else {
		log.Printf("Processing IssueData:: %s", cloudEvent.ObjectKey)
	}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}

	data.Description = util.ChangeString(data.Description)
	data.Subject = util.ChangeString(data.Subject)

	err = util.ChangeTime(
		&data.ClosedOn,
		&data.DueDate,
		&data.CreatedOn,
		&data.UpdatedOn,
		&data.StartDate,
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
	_, err = client.Collection("raw-issues").Doc(strconv.Itoa(data.ID)).Set(ctx, firestoreData)

	return err
}
