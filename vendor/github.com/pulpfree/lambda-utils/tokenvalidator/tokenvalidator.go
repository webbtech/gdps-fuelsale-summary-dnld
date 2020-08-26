package tokenvalidator

import (
	"errors"
	"fmt"

	jwkfetch "github.com/Soluto/fetch-jwk"
	jwt "github.com/dgrijalva/jwt-go"
)

// docs for jwkfetch: https://pkg.go.dev/github.com/Soluto/fetch-jwk?tab=doc
// docs for jwt: https://pkg.go.dev/github.com/dgrijalva/jwt-go?tab=doc
// docs for lambda authorizer: https://github.com/awslabs/aws-apigateway-lambda-authorizer-blueprints

// Validate function
// Other than the following: algorithm, expiry, and token structure, the only other thing
// we're checking is the Cognito clientID
func Validate(cognitoClientID, tokenString string) (principalID string, err error) {

	var errStr, clientID, username string

	jwk := jwkfetch.FromIssuerClaim
	token, err := jwt.Parse(tokenString, jwk())
	claims := token.Claims.(jwt.MapClaims)

	clientID, err = setClientID(claims)
	if err != nil {
		return principalID, err
	}
	username, err = setUsername(claims)
	if err != nil {
		return principalID, err
	}

	// oddly, an expired token doesn't get past this check,
	// so neither jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet get trapped
	if token.Valid == false {
		return principalID, errors.New("Invalid token")
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return principalID, errors.New("Invalid token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return principalID, errors.New("Expired token")
		} else {
			return principalID, errors.New("Invalid token with unknown type")
		}
	}

	// Validate the expected alg
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		errStr = fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
		return principalID, errors.New(errStr)
	}

	// Check clientID
	if cognitoClientID != clientID {
		return principalID, errors.New("Invalid client id")
	}

	// now that we've eliminated any errors, create the id
	principalID = fmt.Sprintf("%s|%s", username, clientID)

	return principalID, err
}

func setClientID(claims jwt.MapClaims) (clientID string, err error) {
	if claims["client_id"] != nil {
		return claims["client_id"].(string), err
	} else if claims["aud"] != nil {
		return claims["aud"].(string), err
	} else {
		return clientID, errors.New("Missing either client_id or aud in claims")
	}
}

func setUsername(claims jwt.MapClaims) (clientID string, err error) {
	if claims["cognito:username"] != nil {
		return claims["cognito:username"].(string), err
	} else if claims["username"] != nil {
		return claims["username"].(string), err
	} else {
		return clientID, errors.New("Missing either client:username or username in claims")
	}
}
