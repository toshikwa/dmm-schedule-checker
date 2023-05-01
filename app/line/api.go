package line

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	accessToken = os.Getenv("LINE_ACCESS_TOKEN")
	lineUrl     = "https://notify-api.line.me/api/notify"
)

func SendMessage(message string) error {
	u, err := url.ParseRequestURI(lineUrl)
	if err != nil {
		return fmt.Errorf("failed to parse request uri: %v", err)
	}

	c := &http.Client{}
	form := url.Values{}
	form.Add("message", "\n"+message)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return fmt.Errorf("failed to crate new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message to line: %v", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return nil
}
