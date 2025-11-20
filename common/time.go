package common

import (
	"database/sql/driver"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type MyTime time.Time

func (t *MyTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = MyTime(time.Time{})
		return
	}
	now, _ := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = MyTime(now)
	return
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t MyTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *MyTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = MyTime(tTime)
	return nil
}

func (t MyTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
