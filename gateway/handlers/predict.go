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

func MakePredictHandler(predictorURL url.URL, next http.HandlerFunc) http.HandlerFunc {

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
			}
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Read response body failed, err: %s\n", err)
			} else {
				log.Printf("Read response body succeed. body: %s", string(responseBody))
			}

		}
		next.ServeHTTP(w, r)
	}
}
