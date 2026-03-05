package utils

import (
	"time"
)

type JsonTime time.Time

const layout = "2006-01-02 15:04:05"

func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format(layout)
	return []byte(`"` + formatted + `"`), nil
}

func (t *JsonTime) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(`"`+layout+`"`, string(data))
	if err != nil {
		return err
	}
	*t = JsonTime(parsed)
	return nil
}
