package function

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

//Response response format for slack API
type Response struct {
	FulfillmentText string `json:"fulfillmentText"`
}

type Intent struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type Request struct {
	ResponseID  string      `json:"responseID"`
	QueryResult QueryResult `json:"queryResult"`
}

type QueryResult struct {
	QueryText                 string  `json:"queryText"`
	AllRequiredParamsPresent  bool    `json:"allRequiredParamsPresent"`
	Intent                    Intent  `json:"intent"`
	IntentDetectionConfidence float32 `json:"intentDetectionConfidence"`
	LanguageCode              string  `json:"languageCode"`
	Session                   string  `json:"session"`
}

// Handle a serverless request
func Handle(req []byte) string {
	os.Stderr.WriteString("Requets: " + string(req))

	var requestPayload Request
	err := json.Unmarshal(req, &requestPayload)
	if err != nil {
		os.Stderr.WriteString("Error during parsing request: " + err.Error())
		return "Unable to parse requets"

	}
	gateway_hostname := os.Getenv("gateway_hostname")
	if gateway_hostname == "" {
		gateway_hostname = "gateway"
	}
	var response Response
	var res *http.Response

	if requestPayload.QueryResult.Intent.DisplayName == "dnd_info" {
		res, err = http.Get("http://" + gateway_hostname + ":8080/function/dnd-info")
	} else if requestPayload.QueryResult.Intent.DisplayName == "end_dnd" {
		res, err = http.Get("http://" + gateway_hostname + ":8080/function/end-dnd")
	} else {
		response = Response{FulfillmentText: "Sorry, I can not do that right now, I am still learning."}
	}

	if err != nil {
		os.Stderr.WriteString("Error during making request: " + err.Error())
		response = Response{FulfillmentText: "Soory, Something went wrong."}
	} else {
		if res != nil {
			if res.StatusCode == 200 {
				defer res.Body.Close()
				resBytes, _ := ioutil.ReadAll(res.Body)
				response = Response{FulfillmentText: string(resBytes)}
			} else {
				response = Response{FulfillmentText: "Soory, Something went wrong."}
			}
		}
	}
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes)
}
