package function

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

type Duration struct {
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Parameters struct {
	Duration Duration `json:"duration"`
}

type QueryResult struct {
	QueryText                 string     `json:"queryText"`
	AllRequiredParamsPresent  bool       `json:"allRequiredParamsPresent"`
	Intent                    Intent     `json:"intent"`
	IntentDetectionConfidence float32    `json:"intentDetectionConfidence"`
	LanguageCode              string     `json:"languageCode"`
	Session                   string     `json:"session"`
	Parameters                Parameters `json:"parameters"`
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

	gatewayHostname := os.Getenv("gateway_hostname")
	if gatewayHostname == "" {
		gatewayHostname = "gateway.openfaas"
	}

	var response Response
	var res *http.Response
	intentName := requestPayload.QueryResult.Intent.DisplayName

	if intentName == "dnd_info" {
		res, err = http.Get("http://" + gatewayHostname + ":8080/function/dnd-info")

	} else if intentName == "end_dnd" {
		res, err = http.Get("http://" + gatewayHostname + ":8080/function/end-dnd")

	} else if intentName == "set_snooze" {
		unit := requestPayload.QueryResult.Parameters.Duration.Unit
		duration := int(requestPayload.QueryResult.Parameters.Duration.Amount)

		if unit == "s" {
			duration = int(duration / 60)
		} else if unit == "h" {
			duration = duration * 60
		}
		reader := bytes.NewReader([]byte(strconv.Itoa(duration)))

		res, err = http.Post("http://"+gatewayHostname+":8080/function/set-snooze", "text/plain", reader)

	} else if intentName == "end_snooze" {
		res, err = http.Get("http://" + gatewayHostname + ":8080/function/end-snooze")

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
