package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

const NanoleafPort = 16021

type Source struct {
	IPAddress string `json:"ip_address"`
	APIToken  string `json:"api_token"`
}

type Params struct {
	On         *bool `json:"power"`
	Hue        *int  `json:"hue"`
	Brightness *int  `json:"brightness"`
}

type PutRequest struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type PutRequestResponseVersion struct {
	Ok string `json:"ok"`
}

type PutRequestResponse struct {
	Version PutRequestResponseVersion
}

type NanoleafStateRequestOn struct {
	Value bool `json:"value,omitempty"`
}

type NanoleafStateRequestHue struct {
	Value int `json:"value,omitempty"`
}

type NanoleafStateRequestBrightness struct {
	Value int `json:"value,omitempty"`
}

type NanoleafStateRequest struct {
	On         *NanoleafStateRequestOn         `json:"on,omitempty"`
	Hue        *NanoleafStateRequestHue        `json:"hue,omitempty"`
	Brightness *NanoleafStateRequestBrightness `json:"brightness,omitempty"`
}

func validate(request PutRequest) (PutRequest, error) {
	if request.Source.IPAddress == "" {
		return PutRequest{}, fmt.Errorf("ip address is required")
	}

	if request.Source.APIToken == "" {
		return PutRequest{}, fmt.Errorf("api token is required")
	}

	if request.Params.On == nil {
		true := true
		request.Params.On = &true
	}

	if request.Params.Hue != nil && (*request.Params.Hue < 1 || *request.Params.Hue > 360) {
		return PutRequest{}, fmt.Errorf("invalid hue setting")
	}

	if request.Params.Brightness != nil && (*request.Params.Brightness < 0 || *request.Params.Brightness > 100) {
		return PutRequest{}, fmt.Errorf("invalid brightness setting")
	}

	return request, nil
}

func doRequest(request PutRequest) error {
	c := resty.New()

	b := &NanoleafStateRequest{}

	if request.Params.On != nil {
		b.On = &NanoleafStateRequestOn{Value: *request.Params.On}
	}

	if request.Params.Hue != nil {
		b.Hue = &NanoleafStateRequestHue{Value: *request.Params.Hue}
	}

	if request.Params.Brightness != nil {
		b.Brightness = &NanoleafStateRequestBrightness{Value: *request.Params.Brightness}
	}

	log.Println("sending state request")

	r, err := c.R().
		SetBody(b).
		SetHeader("Content-Type", "application/json").
		Put(fmt.Sprintf("http://%s:%d/api/v1/%s/state", request.Source.IPAddress, NanoleafPort, request.Source.APIToken))

	if err != nil {
		return err
	}

	if r.IsError() {
		return fmt.Errorf("bad http response code: %d", r.StatusCode())
	}

	return nil
}

func main() {
	var putRequest PutRequest
	err := json.NewDecoder(os.Stdin).Decode(&putRequest)
	if err != nil {
		log.Fatal("couldn't decode stdin: ", err)
	}

	putRequest, err = validate(putRequest)
	if err != nil {
		log.Fatal("error validating options:", err)
	}

	err = doRequest(putRequest)

	if err != nil {
		log.Fatal("error whilst executing request:", err)
	}

	err = json.NewEncoder(os.Stdout).Encode(PutRequestResponse{Version: PutRequestResponseVersion{Ok: "true"}})
	if err != nil {
		log.Fatal("couldn't encode response", err)
	}
}
