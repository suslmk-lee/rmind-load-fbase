package util

import (
	"rmind-load-fbase/common"
	"rmind-load-fbase/internal/model"
	"time"
)

// 치환할 값들을 정의한 map
var replacements = map[string]string{
	"h3. ": "- ",
	"h2. ": "- ",
	"h1. ": "- ",
	"h4. ": "- ",
	"h5. ": "- ",
	"h6. ": "- ",
	"\r":   "",
	"> ":   "",
}

func ChangeTimeZone(data model.IssueData) (error, model.IssueData) {
	seoulLocation, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return err, data
	}

	data.DueDate = time.Date(
		data.DueDate.Year(), data.DueDate.Month(), data.DueDate.Day(),
		data.DueDate.Hour(), data.DueDate.Minute(), data.DueDate.Second(), data.DueDate.Nanosecond(), seoulLocation).UTC()
	data.CreatedOn = time.Date(
		data.CreatedOn.Year(), data.CreatedOn.Month(), data.CreatedOn.Day(),
		data.CreatedOn.Hour(), data.CreatedOn.Minute(), data.CreatedOn.Second(), data.CreatedOn.Nanosecond(), seoulLocation).UTC()
	data.UpdatedOn = time.Date(
		data.UpdatedOn.Year(), data.UpdatedOn.Month(), data.UpdatedOn.Day(),
		data.UpdatedOn.Hour(), data.UpdatedOn.Minute(), data.UpdatedOn.Second(), data.UpdatedOn.Nanosecond(), seoulLocation).UTC()
	data.StartDate = time.Date(
		data.StartDate.Year(), data.StartDate.Month(), data.StartDate.Day(),
		data.StartDate.Hour(), data.StartDate.Minute(), data.StartDate.Second(), data.StartDate.Nanosecond(), seoulLocation).UTC()

	return nil, data
}

func ChangeString(data model.IssueData) model.IssueData {
	data.Notes = common.ReplaceOrRemove(data.Notes, replacements)
	data.Subject = common.ReplaceOrRemove(data.Subject, replacements)
	data.Description = common.ReplaceOrRemove(data.Description, replacements)

	return data
}
