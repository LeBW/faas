package handlers

import (
	"bytes"
	"encoding/json"
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
			//predictReq, err := http.NewRequest(http.MethodPost, predictorURL.String(), bytes.NewBuffer(requestBody))
			resp, err := http.Post(predictorURL.String(), "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				log.Printf("Send predict request failed, err: %s\n", err)
			} else {
				log.Printf("Send predict request succeed. Result: %s\n", resp.Body)
			}
		}
		next.ServeHTTP(w, r)
	}
}
