package function

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

//Request struct for request body
type Request struct {
	NumberOfMinutes int `json:"num_minutes"`
}

//Response struct for response body
type Response struct {
	OK              bool `json:"ok"`
	SnoozeEnabled   bool `json:"snooze_enabled"`
	SnoozeEndtime   int  `json:"snooze_endtime"`
	SnoozeRemaining int  `json:"snooze_remaining"`
}

// Handle a serverless request
func Handle(req []byte) string {
	URL := "https://slack.com/api/dnd.setSnooze"
	URL = URL + "?num_minutes=" + string(req)

	slackToken := os.Getenv("slack_token")
	if slackToken == "" {
		os.Stderr.WriteString("Slack token is not set")
		return ""
	}

	client := http.Client{}
	request, _ := http.NewRequest(http.MethodGet, URL, nil)
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

		if resp.SnoozeEnabled == true {
			return "Snooze has been enabled for " + string(req) + " minutes"
		}
	}
	return "Something went wrong"
}
