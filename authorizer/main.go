package main

// Copyright 2015-2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License").
// You may not use this file except in compliance with the License.
// A copy of the License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// This file originally found at: https://github.com/awslabs/aws-apigateway-lambda-authorizer-blueprints/blob/master/blueprints/go/main.go

import (
	"context"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pulpfree/gdps-fs-sum-dwnld/config"
	"github.com/pulpfree/lambda-utils/tokenvalidator"
)

var cfg *config.Config

func init() {
	cfg = &config.Config{}
	err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// log.Println("Client token: " + event.AuthorizationToken)
	// log.Println("Method ARN: " + event.MethodArn)

	principalID, err := tokenvalidator.Validate(cfg.CognitoClientID, event.AuthorizationToken)
	log.Println("Validating token")
	if err != nil {
		// log.Println("Error in token validation: " + err.Error())
		log.Errorf("Error in token validation: %+v" + err.Error())
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	// validate the incoming token
	// and produce the principal user identifier associated with the token

	// this could be accomplished in a number of ways:
	// 1. Call out to OAuth provider
	// 2. Decode a JWT token inline
	// 3. Lookup in a self-managed DB
	// principalID := "user|a1b2c3d4"
	// you can send a 401 Unauthorized response to the client by failing like so:
	// return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")

	// if the token is valid, a policy must be generated which will allow or deny access to the client

	// if access is denied, the client will receive a 403 Access Denied response
	// if access is allowed, API Gateway will proceed with the backend integration configured on the method that was called

	// this function must generate a policy that is associated with the recognized principal user identifier.
	// depending on your use case, you might store policies in a DB, or generate them on the fly

	// keep in mind, the policy is cached for 5 minutes by default (TTL is configurable in the authorizer)
	// and will apply to subsequent calls to any method/resource in the RestApi
	// made with the same token

	//the example policy below denies access to all resources in the RestApi
	tmp := strings.Split(event.MethodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")
	awsAccountID := tmp[4]

	resp := NewAuthorizerResponse(principalID, awsAccountID)
	resp.Region = tmp[3]
	resp.APIID = apiGatewayArnTmp[0]
	resp.Stage = apiGatewayArnTmp[1]
	// resp.DenyAllMethods()
	// resp.AllowMethod(Get, "/pets/*")
	resp.AllowAllMethods()

	// new! -- add additional key-value pairs associated with the authenticated principal
	// these are made available by APIGW like so: $context.authorizer.<key>
	// additional context is cached
	/* resp.Context = map[string]interface{}{
		"stringKey":  "stringval",
		"numberKey":  123,
		"booleanKey": true,
	} */

	return resp.APIGatewayCustomAuthorizerResponse, nil
}

func main() {
	// lambda.Start(thundra.Wrap(handleRequest))
	lambda.Start(handleRequest)
}

// HTTPVerb type
type HTTPVerb int

// Verb constants
const (
	Get HTTPVerb = iota
	Post
	Put
	Delete
	Patch
	Head
	Options
	All
)

func (hv HTTPVerb) String() string {
	switch hv {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	case Patch:
		return "PATCH"
	case Head:
		return "HEAD"
	case Options:
		return "OPTIONS"
	case All:
		return "*"
	}
	return ""
}

// Effect type
type Effect int

// Action constant
const (
	Allow Effect = iota
	Deny
)

func (e Effect) String() string {
	switch e {
	case Allow:
		return "Allow"
	case Deny:
		return "Deny"
	}
	return ""
}

// AuthorizerResponse struct
type AuthorizerResponse struct {
	events.APIGatewayCustomAuthorizerResponse

	// The region where the API is deployed. By default this is set to '*'
	Region string

	// The AWS account id the policy will be generated for. This is used to create the method ARNs.
	AccountID string

	// The API Gateway API id. By default this is set to '*'
	APIID string

	// The name of the stage used in the policy. By default this is set to '*'
	Stage string
}

// NewAuthorizerResponse function
func NewAuthorizerResponse(principalID string, AccountID string) *AuthorizerResponse {
	return &AuthorizerResponse{
		APIGatewayCustomAuthorizerResponse: events.APIGatewayCustomAuthorizerResponse{
			PrincipalID: principalID,
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
			},
		},
		Region:    "*",
		AccountID: AccountID,
		APIID:     "*",
		Stage:     "*",
	}
}

func (r *AuthorizerResponse) addMethod(effect Effect, verb HTTPVerb, resource string) {
	resourceArn := "arn:aws:execute-api:" +
		r.Region + ":" +
		r.AccountID + ":" +
		r.APIID + "/" +
		r.Stage + "/" +
		verb.String() + "/" +
		strings.TrimLeft(resource, "/")

	s := events.IAMPolicyStatement{
		Effect:   effect.String(),
		Action:   []string{"execute-api:Invoke"},
		Resource: []string{resourceArn},
	}

	r.PolicyDocument.Statement = append(r.PolicyDocument.Statement, s)
}

// AllowAllMethods method
func (r *AuthorizerResponse) AllowAllMethods() {
	r.addMethod(Allow, All, "*")
}

// DenyAllMethods method
func (r *AuthorizerResponse) DenyAllMethods() {
	r.addMethod(Deny, All, "*")
}

// AllowMethod method
func (r *AuthorizerResponse) AllowMethod(verb HTTPVerb, resource string) {
	r.addMethod(Allow, verb, resource)
}

// DenyMethod method
func (r *AuthorizerResponse) DenyMethod(verb HTTPVerb, resource string) {
	r.addMethod(Deny, verb, resource)
}
