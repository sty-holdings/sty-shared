// Package sty_shared
/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS

    {Other catagories of restrictions}
    * {List of restrictions for the catagory

NOTES:
    {Enter any additional notes that you believe will help the next developer.}

COPYRIGHT:
	Copyright 2022
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.

*/
package sty_shared

import (
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang-jwt/jwt/v5"
)

// Config ...
type Config struct {
	CognitoRegion     string
	CognitoUserPoolID string
}

type KeySet struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

type AWSConfig struct {
	ClientId   string
	Region     string
	UserPoolId string
}

type AWSHelper struct {
	InfoFQN    string
	KeySetURL  string
	KeySet     KeySet
	AWSConfig  AWSConfig
	SessionPtr *awsSession.Session
	tokenType  string
}

type Claims struct {
	AtHash              string `json:"at_hash"`
	AuthTime            int    `json:"auth_time"`
	CognitoUsername     string `json:"cognito:username"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumber         string `json:"phone_number"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	TokenUse            string `json:"token_use"`
	UserName            string `json:"username"`
	jwt.RegisteredClaims
}

var (
	awsConfig = AWSConfig{
		ClientId:   "4i4onptb55891872nfc00bk30a",
		Region:     "us-west-2",
		UserPoolId: "us-west-2_lvAuSOXGf",
	}
)
