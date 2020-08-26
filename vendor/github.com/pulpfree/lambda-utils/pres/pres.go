package pres

// This package provides a struct and wrapper function for
// the aws-lambda-go/events.APIGatewayProxyResponse function

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pulpfree/lambda-utils/pkgerrors"
)

// Response struct
type Response struct {
	Code      int         `json:"code"`      // HTTP status code
	Data      interface{} `json:"data"`      // Data payload
	Message   string      `json:"message"`   // Error or status message
	Status    string      `json:"status"`    // Status code (error|fail|success)
	Timestamp int64       `json:"timestamp"` // Machine-readable UTC timestamp in nanoseconds since EPOCH
}

var stdError *pkgerrors.StdError

// ProxyRes function
func ProxyRes(resp Response, hdrs map[string]string, err error) events.APIGatewayProxyResponse {

	if err != nil {
		resp.Code = 500
		resp.Status = "error"
		log.Error(err)
		// send friendly error to client
		if ok := errors.As(err, &stdError); ok {
			resp.Message = stdError.Msg
		} else {
			resp.Message = err.Error()
		}
	}
	body, _ := json.Marshal(&resp)

	return events.APIGatewayProxyResponse{Body: string(body), Headers: hdrs, StatusCode: resp.Code}
}
