package utils

import (
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
