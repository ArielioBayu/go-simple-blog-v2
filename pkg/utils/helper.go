package utils

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (t JsonTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *JsonTime) Scan(src any) error {
	if src == nil {
		*t = JsonTime(time.Time{})
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		*t = JsonTime(v)
		return nil
	case []byte:
		return t.parseString(string(v))
	case string:
		return t.parseString(v)
	default:
		return fmt.Errorf("cannot scan %T into JsonTime", src)
	}
}

func (t *JsonTime) parseString(str string) error {
	layouts := []string{
		layout,
		time.RFC3339,
		"2006-01-02 15:04:05.000000",
		"2006-01-02",
	}
	for _, l := range layouts {
		if parsed, err := time.Parse(l, str); err == nil {
			*t = JsonTime(parsed)
			return nil
		}
	}
	return fmt.Errorf("cannot parse %q into JsonTime", str)
}

func SetAccessTokenCookie(c *gin.Context, token string) {
	SetCookie(c, "access_token", token, 1*time.Minute)
}

func SetCookie(c *gin.Context, name, token string, maxExpired time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		MaxAge:   int(maxExpired.Seconds()),
		HttpOnly: false,
	}

	http.SetCookie(c.Writer, cookie)
}

func GetCookie() {

}
