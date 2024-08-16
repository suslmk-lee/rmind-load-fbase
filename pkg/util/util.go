package util

import (
	"rmind-load-fbase/common"
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

func ChangeTimeZone(data time.Time) (error, time.Time) {
	seoulLocation, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return err, data
	}

	data = time.Date(
		data.Year(), data.Month(), data.Day(),
		data.Hour(), data.Minute(), data.Second(), data.Nanosecond(), seoulLocation).UTC()

	return nil, data
}

func ChangeTime(data ...*time.Time) error {

	for _, data := range data {
		err, changeTime := ChangeTimeZone(*data)
		if err != nil {
			return err
		}
		*data = changeTime
	}

	return nil

}

func ChangeString(chgStr string) string {
	chgStr = common.ReplaceOrRemove(chgStr, replacements)
	return chgStr
}
