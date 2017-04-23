package main

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/tobyjsullivan/moneypenny/updates"
    "errors"
)

const (
    IntentGeneralUpdate = "GeneralUpdate"
)

func alexaRequestHandler(w http.ResponseWriter, r *http.Request) {
	println("Request received:")
	
	println(fmt.Sprintf("%s %s", r.Method, r.URL.String()))	

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
    reqBody := buf.String()
	println(reqBody)

    reqPayload := &alexaRequestPayload{}
    err := json.Unmarshal([]byte(reqBody),reqPayload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    out, err := routeRequest(reqPayload.Request)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	resp := &alexaResponsePayload{
		Version: "1.0",
		Response: &response {
			OutputSpeech: &outputSpeech {
				Type: "PlainText",
				Text: out,
			},
		},
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func routeRequest(req *request) (string, error) {
    switch req.Intent.Name {
    case IntentGeneralUpdate:
        return updates.BuildResponse(), nil
    default:
        return "", errors.New("Unsupported intent.")
    }
}

type alexaRequestPayload struct {
    Request *request `json:"request"`
}

type request struct {
    Intent *intent `json:"intent"`
}

type intent struct {
    Name string `json:"name"`
}

type alexaResponsePayload struct {
	Version string `json:"version"`
	Response *response `json:"response"`
}

type response struct{
	OutputSpeech *outputSpeech `json:"outputSpeech"`
}

type outputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}