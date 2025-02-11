package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func MakePredictHandler(predictorURL url.URL, scheduler FunctionScheduler, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("In MakePredictHandler. r.Host: %s, r.RemoteAddress: %s\n", r.Host, r.RemoteAddr)
		if r.Host != "gateway.openfaas:8080" {
			splits := strings.Split(r.URL.Path, "/")
			functionName := splits[len(splits)-1]

			reqParams := make(map[string]string)
			for key, values := range r.URL.Query() {
				reqParams[key] = values[0]
			}
			reqParamsJson, err := json.Marshal(reqParams)
			if err != nil {
				log.Printf("Parse reqParams to json failedl. err: %s\n", err)
			}

			values := map[string]string{
				"function_name": functionName,
				"req_params":    string(reqParamsJson),
			}
			requestBody, err := json.Marshal(values)
			if err != nil {
				log.Printf("Parse map to request body failed, err: %s\n", err)
			}
			log.Printf("[MakePredictHandler] request body: %s", requestBody)
			finalURL := predictorURL.String() + "/predict"
			log.Printf("[MakePredictHandler] finalURL: %s\n", finalURL)
			resp, err := http.Post(finalURL, "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				log.Printf("Send predict request failed, err: %s\n", err)
			} else if resp.StatusCode != 200 {
				log.Printf("Send predict request failed. status code: %d\n", resp.StatusCode)
			} else {
				log.Printf("Send predict request succeed\n")
				responseBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Read response body failed, err: %s\n", err)
				} else {
					log.Printf("Read response body succeed. body: %s", string(responseBody))
					var pr PredictResponse
					err := json.Unmarshal(responseBody, &pr)
					if err != nil {
						log.Printf("Unmarshal responsebody to json failed, %s\n", err)
					}
					scheduler.AddPredictions(pr.Data)
				}
			}

		}
		next.ServeHTTP(w, r)
	}
}

type Prediction struct {
	FunctionName string  `json:"function_name"`
	PredictTime  float64 `json:"predict_time"`
	Probability  float64 `json:"probability"`
	ResponseTime float64 `json:"response_time"`
}

type PredictResponse struct {
	Data []Prediction `json:"data"`
}
