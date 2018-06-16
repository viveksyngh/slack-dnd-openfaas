package function

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	OK bool `json:"ok"`
}

// Handle a serverless request
func Handle(req []byte) string {
	URL := "https://slack.com/api/dnd.endDnd"
	slackToken := os.Getenv("slack_token")
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
			return "DnD has been deactivated."
		}
		return "Failed to decativate DND, Something went wrong."
	}

	return "Failed to decativate DND, Something went wrong."
}
