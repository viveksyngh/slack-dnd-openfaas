package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//Response struct response for
type Response struct {
	OK               bool `json:"ok"`
	DnDEnabled       bool `json:"dnd_enabled"`
	NextDnDStartTime int  `json:"next_dnd_start_ts"`
	NextDnDEndTime   int  `json:"next_dnd_end_ts"`
	SnoozeEnabled    bool `json:"snooze_enabled"`
}

// Handle a serverless request
func Handle(req []byte) string {
	URL := "https://slack.com/api/dnd.endSnooze"
	slackToken := os.Getenv("slack_token")
	if slackToken == "" {
		slackTokenBytes, err := ioutil.ReadFile("/var/openfaas/secrets/slack-token")
		slackToken = string(slackTokenBytes)
		if err != nil {
			return fmt.Sprintf("Unable to read secret file")
		}
	}

	if slackToken == "" {
		os.Stderr.WriteString("Slack token is not set")
		return ""
	}
	client := http.Client{}

	request, _ := http.NewRequest(http.MethodPost, URL, nil)
	request.Header.Add("Authorization", "Bearer "+slackToken)
	request.Header.Add("Content-type", "application/x-www-form-urlencoded")

	res, err := client.Do(request)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return "Something went wrong: " + err.Error()
	}

	if res.StatusCode == 200 {
		defer res.Body.Close()

		responseBytes, _ := ioutil.ReadAll(res.Body)
		os.Stderr.WriteString("Response: " + string(responseBytes))
		var resp Response
		err = json.Unmarshal(responseBytes, &resp)

		if err != nil {
			os.Stderr.WriteString(err.Error())
			return "Something went wrong"
		}

		if resp.OK == true {
			return "Snooze has been disabled."
		}
		return "Failed to end snooze, Something went wrong."
	}

	return "Failed to end snooze, Something went wrong."
}
