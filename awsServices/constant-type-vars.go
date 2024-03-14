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
	"context"
	"math/big"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// https://github.com/aws/amazon-cognito-identity-js/blob/master/src/AuthenticationHelper.js#L22
	nHex = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"
	// https://github.com/aws/amazon-cognito-identity-js/blob/master/src/AuthenticationHelper.js#L49
	gHex     = "2"
	infoBits = "Caldera Derived Key"
)

type AWSSession struct {
	baseConfig       aws.Config
	clientConfig     clientConfig
	identityPoolInfo cognitoIdentityInfo
	keyInfo          cognitoKeyInfo
	styBaseConfig    awsConfig
	tokens           cognitoTokens
	username         string
}

type awsConfig struct {
	clientId       string
	identityPoolId string
	region         string
	userPoolId     string
}

type clientConfig struct {
	clientId string
	username string
}

// ToDo Is this needed
type cognitoClaims struct {
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

type cognitoIdentityInfo struct {
	credentials aws.Credentials
	IdentityId  string
}

type cognitoKeyInfo struct {
	keySet    keySet
	keySetURL string
}

// CognitoLogin handles SRP authentication with AWS Cognito
type cognitoLogin struct {
	username     string
	password     *string
	userPoolId   string
	userPoolName string
	clientId     string
	clientSecret *string
	bigN         *big.Int
	g            *big.Int
	k            *big.Int
	a            *big.Int
	bigA         *big.Int
}

type cognitoTokens struct {
	access  string
	id      string
	refresh string
}

type keySet struct {
	Keys []key `json:"keys"`
}

type key struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
}

var (
	styConfigLocal = awsConfig{
		clientId:       "677jaef1i0cri2hpbtvfce4152",
		identityPoolId: "",
		region:         "us-west-2",
		userPoolId:     "us-west-2_lvAuSOXGf",
	}
	styConfigDevelopment = awsConfig{
		clientId:       "5nfnlbaoiprg5q5n7jvd2lvm0d",
		identityPoolId: "",
		region:         "us-west-2",
		userPoolId:     "us-west-2_d0U66vAT1",
	}
	styConfigProduction = awsConfig{
		clientId:       "4i4onptb55891872nfc00bk30a",
		identityPoolId: "us-west-2:973d66b8-dece-4315-8ab3-58ad924b357b",
		region:         "us-west-2",
		userPoolId:     "us-west-2_lvAuSOXGf",
	}
)

var (
	awsCTXToDo       = context.TODO()
	awsCTXBackground = context.Background()
	awsTrue          = true
	awsTruePtr       = &awsTrue
)
