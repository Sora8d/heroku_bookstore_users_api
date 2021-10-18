package date

import "time"

const apiDateLayout = "2006-01-02 15:04:05"

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	date := GetNow()
	return date.Format(apiDateLayout) + " UTC"
}
