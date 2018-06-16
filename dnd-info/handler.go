package function

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

//Response response format for slack API
type Response struct {
	OK            bool `json:"ok"`
	DNDEnabled    bool `json:"dnd_enabled"`
	SnoozeEnabled bool `json:"snooze_enabled"`
}

// Handle a serverless request
func Handle(req []byte) string {
	URL := "https://slack.com/api/dnd.info"
	slackToken := os.Getenv("slack_token")
	if slackToken == "" {
		os.Stderr.WriteString("Slack token is not set")
		return ""
	}
	URL = URL + "?token=" + slackToken

	res, err := http.Get(URL)
	if err != nil {
		os.Stderr.WriteString("Unable to get DnD info")
		return "Failed to DnD info"
	}
	if res.StatusCode == 200 {
		defer res.Body.Close()
		responseBytes, _ := ioutil.ReadAll(res.Body)
		var resp Response
		err = json.Unmarshal(responseBytes, &resp)
		if resp.DNDEnabled == true {
			return "You have DND enabled."
		}
		return "You do not have DnD enabled."
	}
	return "Something went wrong."
}
