package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//Response response format for slack API
type Response struct {
	OK              bool `json:"ok"`
	DNDEnabled      bool `json:"dnd_enabled"`
	SnoozeEnabled   bool `json:"snooze_enabled"`
	SnoozeRemaining int  `json:"snooze_remaining"`
}

// Handle a serverless request
func Handle(req []byte) string {
	URL := "https://slack.com/api/dnd.info"
	slackToken := os.Getenv("slack_token")
	if slackToken == "" {
		os.Stderr.WriteString("Slack token is not set")
		return ""
	}
	os.Stderr.WriteString(string(req))
	URL = URL + "?token=" + slackToken

	res, err := http.Get(URL)
	if err != nil {
		os.Stderr.WriteString("Unable to get DnD info")
		return "Failed to DnD info"
	}
	if res.StatusCode == 200 {
		defer res.Body.Close()
		responseBytes, _ := ioutil.ReadAll(res.Body)
		os.Stderr.WriteString(string(responseBytes))
		var resp Response
		err = json.Unmarshal(responseBytes, &resp)
		var responseMessage string
		if resp.DNDEnabled == true {
			responseMessage = "You have DND enabled."
		} else {
			responseMessage = "You do not have DND enabled."
		}

		if resp.SnoozeEnabled == true {
			minutes := resp.SnoozeRemaining / 60
			seconds := resp.SnoozeRemaining % 60
			responseMessage = responseMessage + fmt.Sprintf("You also have active snooze notifications which will end in %d minutes and %d seconds.", minutes, seconds)
		} else {
			responseMessage = responseMessage + "You don't have active snooze notifications."
		}
		return responseMessage
	}
	return "Something went wrong."
}
