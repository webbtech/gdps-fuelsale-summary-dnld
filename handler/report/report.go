package main

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/pulpfree/gdps-fs-sum-dwnld/config"
	"github.com/pulpfree/gdps-fs-sum-dwnld/model"
	"github.com/pulpfree/gdps-fs-sum-dwnld/report"
	"github.com/pulpfree/gdps-fs-sum-dwnld/validate"
	"github.com/pulpfree/lambda-utils/pres"
)

var cfg *config.Config

func init() {
	cfg = &config.Config{}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

// SignedURL struct
type SignedURL struct {
	URL string `json:"url"`
}

// HandleRequest function
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	hdrs["Access-Control-Allow-Origin"] = "*"
	hdrs["Access-Control-Allow-Methods"] = "GET,OPTIONS,POST,PUT"
	hdrs["Access-Control-Allow-Headers"] = "Authorization,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{Body: string("null"), Headers: hdrs, StatusCode: 200}, nil
	}

	t := time.Now()

	// If this is a ping test, intercept and return
	if req.HTTPMethod == "GET" {
		log.Info("Ping test in handleRequest")
		return pres.ProxyRes(pres.Response{
			Code:      200,
			Data:      "pong",
			Status:    "success",
			Timestamp: t.Unix(),
		}, hdrs, nil), nil
	}

	// Set and validate request params
	var r *model.RequestInput
	json.Unmarshal([]byte(req.Body), &r)
	reqVars, err := validate.RequestInput(r)
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	// Process request
	report, err := report.New(reqVars, cfg, req.Headers["Authorization"])
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	// var url string
	err = report.Create()
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}

	url, err := report.CreateSignedURL()
	if err != nil {
		return pres.ProxyRes(pres.Response{
			Timestamp: t.Unix(),
		}, hdrs, err), nil
	}
	log.Infof("signed url created %s", url)

	return pres.ProxyRes(pres.Response{
		Code:      201,
		Data:      SignedURL{URL: url},
		Status:    "success",
		Timestamp: t.Unix(),
	}, hdrs, nil), nil
}

func main() {
	lambda.Start(HandleRequest)
}
