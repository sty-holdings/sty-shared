// Package sty_shared
/*
This is the STY-Holdings shared services

NOTES:

	None

COPYRIGHT & WARRANTY:

	Copyright (c) 2022 STY-Holdings, inc
	All rights reserved.

	This software is the confidential and proprietary information of STY-Holdings, Inc.
	Use is subject to license terms.

	Unauthorized copying of this file, via any medium is strictly prohibited.

	Proprietary and confidential

	Written by <Replace with FULL_NAME> / syacko
	STY-Holdings, Inc.
	support@sty-holdings.com
	www.sty-holdings.com

	01-2024
	USA

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package sty_shared

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	awsCred "github.com/aws/aws-sdk-go-v2/credentials"
	awsCI "github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	awsCIP "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	awsCT "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	awsSSM "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/golang-jwt/jwt/v5"
	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// GetClientId returns the client ID stored in the sessionPtr.
//
// Customer Messages: None
// Errors: None
// Verifications: None
func GetClientId(sessionPtr *AWSSession) string {

	return sessionPtr.clientConfig.clientId
}

// GetIdentityCredentials - will return AWS temporary credentials.
// The variables 'identityId' is option and are only used when sessionPtr values are empty.
// The variable 'identityIdCredentials' is only needed if sessionPtr is nil.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func GetIdentityCredentials(
	sessionPtr *AWSSession,
	identityId string,

) (
	identityIdCredentials aws.Credentials,
	errorInfo pi.ErrorInfo,
) {

	var (
		tClientPtr                 *awsCI.Client
		tGetIdentityCredentialsPtr *awsCI.GetCredentialsForIdentityOutput
		tIdentityId                string
		tLogins                    = make(map[string]string)
	)

	if sessionPtr.identityPoolInfo.identityId == ctv.VAL_EMPTY {
		if identityId == ctv.VAL_EMPTY {
			errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v or AWSSession Access Token", ctv.TXT_MISSING_PARAMETER, ctv.FN_AWS_IDENTITY_ID))
			return
		}
		tIdentityId = identityId
	} else {
		tIdentityId = sessionPtr.identityPoolInfo.identityId
	}

	if tClientPtr = awsCI.NewFromConfig(sessionPtr.baseConfig); tClientPtr == nil {
		errorInfo = pi.NewErrorInfo(pi.ErrServiceFailedAWS, fmt.Sprintf("%v%v", ctv.TXT_SERVICE, ctv.TXT_AWS_STS))
		return
	}

	tLogins[fmt.Sprintf("cognito-idp.%v.amazonaws.com/%v", sessionPtr.styBaseConfig.region, sessionPtr.styBaseConfig.userPoolId)] = sessionPtr.tokens.id
	if tGetIdentityCredentialsPtr, errorInfo.Error = tClientPtr.GetCredentialsForIdentity(
		awsCTXToDo, &awsCI.GetCredentialsForIdentityInput{
			IdentityId: aws.String(tIdentityId),
			Logins:     tLogins,
		},
	); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, ctv.VAL_EMPTY)
		return
	}

	sessionPtr.identityPoolInfo.identityId = *tGetIdentityCredentialsPtr.IdentityId
	sessionPtr.identityPoolInfo.credentials.AccessKeyID = *tGetIdentityCredentialsPtr.Credentials.AccessKeyId
	sessionPtr.identityPoolInfo.credentials.SessionToken = *tGetIdentityCredentialsPtr.Credentials.SessionToken
	sessionPtr.identityPoolInfo.credentials.SecretAccessKey = *tGetIdentityCredentialsPtr.Credentials.SecretKey
	sessionPtr.identityPoolInfo.credentials.Expires = *tGetIdentityCredentialsPtr.Credentials.Expiration
	identityIdCredentials = sessionPtr.identityPoolInfo.credentials

	return
}

// GetId - will return System Manager parameters. WithDecryption is assumed.
// The variables 'region' and 'userPoolId' are option and are only used when sessionPtr values are empty.
// The variable 'identityId' is only needed if sessionPtr is nil.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func GetId(
	sessionPtr *AWSSession,
	region, userPoolId string,
) (
	identityId string,
	errorInfo pi.ErrorInfo,
) {

	var (
		tClientPtr      *awsCI.Client
		tGetIdOutputPtr *awsCI.GetIdOutput
		tLogins         = make(map[string]string)
		tRegion         string
		tUserPoolId     string
	)

	if sessionPtr.styBaseConfig.region == ctv.VAL_EMPTY {
		if userPoolId == ctv.VAL_EMPTY {
			errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v or AWSSession Access Token", ctv.TXT_MISSING_PARAMETER, ctv.FN_AWS_REGION))
			return
		}
		tRegion = region
	} else {
		tRegion = sessionPtr.styBaseConfig.region
	}
	if sessionPtr.styBaseConfig.userPoolId == ctv.VAL_EMPTY {
		if userPoolId == ctv.VAL_EMPTY {
			errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v or AWSSession Access Token", ctv.TXT_MISSING_PARAMETER, ctv.FN_USERPOOL_ID))
			return
		}
		tUserPoolId = userPoolId
	} else {
		tUserPoolId = sessionPtr.styBaseConfig.userPoolId
	}

	if tClientPtr = awsCI.NewFromConfig(sessionPtr.baseConfig); tClientPtr == nil {
		errorInfo = pi.NewErrorInfo(pi.ErrServiceFailedAWS, fmt.Sprintf("%v%v", ctv.TXT_SERVICE, ctv.TXT_AWS_COGNITO))
	}

	tLogins[fmt.Sprintf("cognito-idp.%v.amazonaws.com/%v", tRegion, tUserPoolId)] = sessionPtr.tokens.id
	if tGetIdOutputPtr, errorInfo.Error = tClientPtr.GetId(
		awsCTXToDo, &awsCI.GetIdInput{
			IdentityPoolId: aws.String(sessionPtr.styBaseConfig.identityPoolId),
			Logins:         tLogins,
		},
	); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, ctv.VAL_EMPTY)
		return
	}

	sessionPtr.identityPoolInfo.identityId = *tGetIdOutputPtr.IdentityId
	identityId = sessionPtr.identityPoolInfo.identityId

	return
}

// GetParameters - will return System Manager parameters. WithDecryption is assumed.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func GetParameters(
	sessionPtr *AWSSession,
	ssmParameters ...string,
) (
	parametersOutput awsSSM.GetParametersOutput,
	errorInfo pi.ErrorInfo,
) {

	var (
		tClientPtr           *awsSSM.Client
		tParametersOutputPtr *awsSSM.GetParametersOutput
	)

	if sessionPtr.styBaseConfig.userPoolId == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_USERPOOL_ID, sessionPtr.styBaseConfig.userPoolId))
		return
	}

	if len(ssmParameters) == ctv.VAL_ZERO {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_MISSING_PARAMETER, ssmParameters))
		return
	}

	if _, errorInfo = GetId(sessionPtr, ctv.VAL_EMPTY, ctv.VAL_EMPTY); errorInfo.Error != nil {
		return
	}

	if _, errorInfo = GetIdentityCredentials(sessionPtr, ctv.VAL_EMPTY); errorInfo.Error != nil {
		return
	}

	sessionPtr.baseConfig.Credentials = awsCred.StaticCredentialsProvider{Value: aws.Credentials{
		AccessKeyID:     sessionPtr.identityPoolInfo.credentials.AccessKeyID,
		SecretAccessKey: sessionPtr.identityPoolInfo.credentials.SecretAccessKey,
		SessionToken:    sessionPtr.identityPoolInfo.credentials.SessionToken,
		Source:          "",
		CanExpire:       false,
		Expires:         time.Time{},
	}}

	if tClientPtr = awsSSM.NewFromConfig(sessionPtr.baseConfig); tClientPtr == nil {
		errorInfo = pi.NewErrorInfo(pi.ErrServiceFailedAWS, fmt.Sprintf("%v%v", ctv.TXT_SERVICE, ctv.TXT_AWS_SYSTEM_MANAGER))
	}

	if tParametersOutputPtr, errorInfo.Error = tClientPtr.GetParameters(
		awsCTXToDo, &awsSSM.GetParametersInput{
			Names:          ssmParameters,
			WithDecryption: awsTruePtr,
		},
	); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, ctv.VAL_EMPTY)
		return
	}

	parametersOutput = *tParametersOutputPtr

	return
}

// Login - authenticates the user with the provided login type, username, and password. It
// stores the tokens in the sessionPtr and any parameters returned from GetParameters
//
// Customer Messages: None
// Errors: None
// Verifications: None
func Login(
	loginType, username string,
	password *string,
	sessionPtr *AWSSession,
) (
	errorInfo pi.ErrorInfo,
) {

	var (
		tClientPtr                    *awsCIP.Client
		cognitoLoginPtr               *cognitoLogin
		tInitiateAuthOutputPtr        *awsCIP.InitiateAuthOutput
		tRespondToAuthChallengeOutput *awsCIP.RespondToAuthChallengeOutput
		tTokens                       = make(map[string]string)
	)

	if loginType == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_LOGIN_TYPE, loginType))
		return
	}
	if username == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_USERNAME, username))
		return
	}
	if password == nil {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_PASSWORD, ctv.TXT_PROTECTED))
		return
	}
	if sessionPtr.styBaseConfig.userPoolId == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_USERPOOL_ID, sessionPtr.styBaseConfig.userPoolId))
		return
	}

	if cognitoLoginPtr, errorInfo = NewCognitoLogin(username, sessionPtr.styBaseConfig.userPoolId, sessionPtr.styBaseConfig.clientId, password, nil); errorInfo.Error != nil {
		pi.PrintErrorInfo(errorInfo)
	}

	if tClientPtr = awsCIP.NewFromConfig(sessionPtr.baseConfig); tClientPtr == nil {
		errorInfo = pi.NewErrorInfo(pi.ErrServiceFailedAWS, fmt.Sprintf("%v%v", ctv.TXT_SERVICE, ctv.TXT_AWS_COGNITO))
	}

	// initiate auth
	if tInitiateAuthOutputPtr, errorInfo.Error = tClientPtr.InitiateAuth(
		context.Background(), &awsCIP.InitiateAuthInput{
			AuthFlow:       awsCT.AuthFlowType(loginType),
			ClientId:       aws.String(cognitoLoginPtr.GetClientId()),
			AuthParameters: cognitoLoginPtr.GetAuthParams(),
		},
	); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, ctv.VAL_EMPTY)
		return
	}

	tTokens = make(map[string]string) // This is used for either awsCT.AuthFlowTypeUserPasswordAuth or awsCT.AuthFlowTypeUserSrpAuth
	if loginType == string(awsCT.AuthFlowTypeUserPasswordAuth) {
		tTokens["access"] = *tInitiateAuthOutputPtr.AuthenticationResult.AccessToken
		tTokens["id"] = *tInitiateAuthOutputPtr.AuthenticationResult.IdToken
		tTokens["refresh"] = *tInitiateAuthOutputPtr.AuthenticationResult.RefreshToken
	}

	// respond to password verifier challenge
	if tInitiateAuthOutputPtr.ChallengeName == awsCT.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := cognitoLoginPtr.PasswordVerifierChallenge(tInitiateAuthOutputPtr.ChallengeParameters, time.Now())
		if tRespondToAuthChallengeOutput, errorInfo.Error = tClientPtr.RespondToAuthChallenge(
			context.Background(), &awsCIP.RespondToAuthChallengeInput{
				ChallengeName:      awsCT.ChallengeNameTypePasswordVerifier,
				ChallengeResponses: challengeResponses,
				ClientId:           aws.String(cognitoLoginPtr.GetClientId()),
			},
		); errorInfo.Error != nil {
			errorInfo = pi.NewErrorInfo(errorInfo.Error, ctv.VAL_EMPTY)
			return
		}
		tTokens["access"] = *tRespondToAuthChallengeOutput.AuthenticationResult.AccessToken
		tTokens["id"] = *tRespondToAuthChallengeOutput.AuthenticationResult.IdToken
		tTokens["refresh"] = *tRespondToAuthChallengeOutput.AuthenticationResult.RefreshToken
	}

	sessionPtr.clientConfig.username = username
	sessionPtr.tokens.access = tTokens["access"]
	sessionPtr.tokens.id = tTokens["id"]
	sessionPtr.tokens.refresh = tTokens["refresh"]

	return
}

// // ConfirmUser - mark the AWS user as confirmed
// func (a *AWSHelper) ConfirmUser(userName string) (errorInfo pi.ErrorInfo) {
//
// 	var (
// 		tAdminConfirmSignUpInput    cognito.AdminConfirmSignUpInput
// 		tawsCIPPtr *cognito.awsCIP
// 		tFunction, _, _, _          = runtime.Caller(0)
// 		tFunctionName               = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	if userName == ctv.VAL_EMPTY {
// 		errorInfo.Error = pi.ErrRequiredArgumentMissing
// 		log.Println(errorInfo.Error)
// 	} else {
// 		tawsCIPPtr = cognito.New(a.SessionPtr)
// 		tAdminConfirmSignUpInput.Username = &userName
// 		tAdminConfirmSignUpInput.UserPoolId = &a.awsCfg.UserPoolId
// 		if _, errorInfo.Error = tawsCIPPtr.AdminConfirmSignUp(&tAdminConfirmSignUpInput); errorInfo.Error != nil {
// 			// If the user is already confirmed, AWS will return an error, and do not care about this error.
// 			if strings.Contains(errorInfo.Error.Error(), ctv.STATUS_CONFIRMED) {
// 				errorInfo.Error = nil
// 			} else {
// 				if strings.Contains(errorInfo.Error.Error(), pi.USER_DOES_NOT_EXIST) {
// 					errorInfo.Error = pi.ErrUserMissing
// 				}
// 			}
// 		}
// 	}
//
// 	return
// }

// GetRequestorEmailPhoneFromIdTokenClaims - will validate the AWS Id JWT, check to make sure the email has been verified, and return the requestor id, email, and phone number.
// func (a *AWSHelper) GetRequestorEmailPhoneFromIdTokenClaims(
// 	firestoreClientPtr *firestore.Client,
// 	token string,
// ) (
// 	requestorId, email, phoneNumber string,
// 	errorInfo pi.ErrorInfo,
// ) {
//
// 	var (
// 		tClaimsPtr         *Claims
// 		tFunction, _, _, _ = runtime.Caller(0)
// 		tFunctionName      = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	if token == ctv.VAL_EMPTY {
// 		errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! Token: '%v'", token))
// 		log.Println(errorInfo.Error)
// 	} else {
// 		if tClaimsPtr, errorInfo = getTokenClaims(a, ctv.TOKEN_TYPE_ID, token); errorInfo.Error == nil {
// 			if isTokenValid(firestoreClientPtr, a, ctv.TOKEN_TYPE_ID, token) {
// 				requestorId = tClaimsPtr.Subject
// 				email = tClaimsPtr.Email
// 				phoneNumber = tClaimsPtr.PhoneNumber
// 			} else {
// 				errorInfo.Error = pi.ErrTokenInvalid
// 				log.Println(errorInfo.Error)
// 			}
// 		}
// 	}
//
// 	return
// }

// GetRequestorFromAccessTokenClaims - will valid the AWS Access JWT, and return the requestor id.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
// func (a *AWSHelper) GetRequestorFromAccessTokenClaims(
// 	firestoreClientPtr *firestore.Client,
// 	token string,
// ) (
// 	requestorId string,
// 	errorInfo pi.ErrorInfo,
// ) {
//
// 	var (
// 		tClaimsPtr         *Claims
// 		tFunction, _, _, _ = runtime.Caller(0)
// 		tFunctionName      = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if token == ctv.TEST_STRING {
// 		requestorId = ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID
// 	} else {
// 		if token == ctv.VAL_EMPTY {
// 			errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! Token: '%v'", token))
// 			log.Println(errorInfo.Error)
// 		} else {
// 			if tClaimsPtr, errorInfo = getTokenClaims(a, ctv.TOKEN_TYPE_ACCESS, token); errorInfo.Error == nil {
// 				if isTokenValid(firestoreClientPtr, a, ctv.TOKEN_TYPE_ACCESS, token) {
// 					requestorId = tClaimsPtr.Subject
// 				} else {
// 					errorInfo.Error = pi.ErrTokenInvalid
// 					log.Println(errorInfo.Error)
// 				}
// 			}
// 		}
// 	}
//
// 	return
// }

// NewAWSConfig - reads the SDKs default external configurations, and populates an AWS Config with the values from the external configurations.
//
//	Customer Messages: None
//	Errors: ErrEnvironmentInvalid, anything awsConfig.LoadDefaultConfig or getPublicKeySet returns,
//	Verifications: None
func NewAWSConfig(environment string) (
	sessionPtr *AWSSession,
	errorInfo pi.ErrorInfo,
) {

	sessionPtr = &AWSSession{} // Initialize session

	switch strings.ToLower(strings.Trim(environment, ctv.SPACES_ONE)) {
	case ctv.ENVIRONMENT_PRODUCTION:
		sessionPtr.styBaseConfig = styConfigProduction
	case ctv.ENVIRONMENT_DEVELOPMENT:
		sessionPtr.styBaseConfig = styConfigDevelopment
	case ctv.ENVIRONMENT_LOCAL:
		sessionPtr.styBaseConfig = styConfigLocal
	default:
		if environment == ctv.VAL_EMPTY {
			errorInfo = pi.NewErrorInfo(pi.ErrEnvironmentInvalid, fmt.Sprintf("%v%v", ctv.TXT_EVIRONMENT, ctv.FN_ENVIRONMENT))
			return
		}
		errorInfo = pi.NewErrorInfo(pi.ErrEnvironmentInvalid, fmt.Sprintf("%v%v", ctv.TXT_EVIRONMENT, environment))
	}

	if sessionPtr.baseConfig, errorInfo.Error = awsCfg.LoadDefaultConfig(awsCTXToDo, awsCfg.WithRegion(sessionPtr.styBaseConfig.region)); errorInfo.
		Error != nil {
		errorInfo = pi.NewErrorInfo(pi.ErrServiceFailedAWS, "Failed to create an AWS Session.")
		return
	}

	sessionPtr.keyInfo.keySetURL = fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", sessionPtr.styBaseConfig.region, sessionPtr.styBaseConfig.userPoolId,
	)
	sessionPtr.keyInfo.keySet, errorInfo = getPublicKeySet(sessionPtr.keyInfo.keySetURL)

	return
}

// ParseAWSJWT - will return the claims, if any, or an err if the AWS JWT is invalid.
// This will parse ID and Access tokens. Refresh token are not support and nothing is returned.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func ParseAWSJWT(
	sessionPtr *AWSSession,
	tokenType, token string,
) (
	claims jwt.Claims,
	tokenValuePtr *jwt.Token,
	errorInfo pi.ErrorInfo,
) {

	if len(sessionPtr.keyInfo.keySet.Keys) == ctv.VAL_ZERO {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, ctv.TXT_KEY_SET_MISSING)
		return
	}
	if tokenType == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_TOKEN_TYPE, ctv.FN_TOKEN_TYPE))
		return
	}
	if token == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_TOKEN, ctv.FN_TOKEN))
		return
	}

	if tokenType == ctv.TOKEN_TYPE_REFRESH {
		return
	}

	for i := 0; i < len(sessionPtr.keyInfo.keySet.Keys); i++ {
		if tokenValuePtr, errorInfo.Error = jwt.ParseWithClaims(
			token, jwt.MapClaims{}, func(token *jwt.Token) (
				key interface{},
				err error,
			) {
				key, err = convertKey(sessionPtr.keyInfo.keySet.Keys[i].E, sessionPtr.keyInfo.keySet.Keys[i].N) // ID
				claims = token.Claims
				return
			},
		); errorInfo.Error != nil {
			fmt.Println(errorInfo.Error)
			if errorInfo.Error.Error() == pi.ErrJWTTokenSignatureInvalid.Error() {
				continue
			} else {
				break
			}
		}
		return // No errors returned from called function
	}

	errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_TOKEN, ctv.FN_TOKEN))

	return
}

// PullCognitoUserInfo - pull user information from the Cognito user pool.
// func (a *AWSHelper) PullCognitoUserInfo(username string) (
// 	userData map[string]interface{},
// 	errorInfo pi.ErrorInfo,
// ) {
//
// 	var (
// 		tAdminGetUserInput          cognito.AdminGetUserInput
// 		tAdminGetUserOutputPtr      *cognito.AdminGetUserOutput
// 		tawsCIPPtr *cognito.awsCIP
// 		tFunction, _, _, _          = runtime.Caller(0)
// 		tFunctionName               = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if username == ctv.VAL_EMPTY {
// 		errorInfo.Error = pi.ErrRequiredArgumentMissing
// 		pi.PrintError(errorInfo)
// 	} else {
// 		tawsCIPPtr = cognito.New(a.SessionPtr)
// 		if tawsCIPPtr == nil {
// 			errorInfo.FileName, errorInfo.ErrorLineNumber = pi.GetFileLineNumber()
// 			errorInfo.Error = pi.ErrServiceFailedAWS
// 			pi.PrintError(errorInfo)
// 		} else {
// 			// Set up the request to get the user
// 			tAdminGetUserInput.UserPoolId = &a.AWSConfig.UserPoolId
// 			tAdminGetUserInput.Username = &username
// 			// Make the request to get the user
// 			if tAdminGetUserOutputPtr, errorInfo.Error = tawsCIPPtr.AdminGetUser(&tAdminGetUserInput); errorInfo.Error == nil {
// 				userData = make(map[string]interface{})
// 				for _, attr := range tAdminGetUserOutputPtr.UserAttributes {
// 					userData[*attr.Name] = *attr.Value
// 				}
// 			} else {
// 				errorInfo.FileName, errorInfo.ErrorLineNumber = pi.GetFileLineNumber()
// 				errorInfo.Error = pi.ErrServiceFailedAWS
// 				pi.PrintError(errorInfo)
// 			}
// 		}
// 	}
//
// 	return
// }

// ValidAWSJWT - will valid the AWS JWT and check to make sure either the phone or email has been verified.
// func (a *AWSHelper) ValidAWSJWT(
// 	firestoreClientPtr *firestore.Client,
// 	tokenType, token string,
// ) (
// 	valid bool,
// 	errorInfo pi.ErrorInfo,
// ) {
//
// 	var (
// 		tFunction, _, _, _ = runtime.Caller(0)
// 		tFunctionName      = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if token == ctv.VAL_EMPTY {
// 		errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! Token: '%v'", token))
// 		log.Println(errorInfo.Error)
// 	} else {
// 		valid = isTokenValid(firestoreClientPtr, a, tokenType, token)
// 	}
//
// 	return
// }

// UpdateAWSEmailVerifyFlag - will update the email_valid field for the user in the Cognito user pool.
// func (a *AWSHelper) UpdateAWSEmailVerifyFlag(username string) (errorInfo pi.ErrorInfo) {
//
// 	var (
// 		tAdminUpdateUserAttributesInput cognito.AdminUpdateUserAttributesInput
// 		tAttributeType                  cognito.AttributeType
// 		tAttributeTypePtrs              []*cognito.AttributeType
// 		tawsCIPPtr     *cognito.awsCIP
// 		tFunction, _, _, _              = runtime.Caller(0)
// 		tFunctionName                   = runtime.FuncForPC(tFunction).Name()
// 		tName                           string
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if username == ctv.VAL_EMPTY {
// 		errorInfo.Error = pi.ErrRequiredArgumentMissing
// 	} else {
// 		tawsCIPPtr = cognito.New(a.SessionPtr)
// 		tName = ctv.FN_EMAIL_VERIFIED // This is required because go doesn't support pointers to ctv.
// 		tAttributeType = cognito.AttributeType{
// 			Name:  &tName,
// 			Value: &tTrueString,
// 		}
// 		tAttributeTypePtrs = append(tAttributeTypePtrs, &tAttributeType)
// 		tAdminUpdateUserAttributesInput.UserAttributes = tAttributeTypePtrs
// 		tAdminUpdateUserAttributesInput.Username = &username
// 		tAdminUpdateUserAttributesInput.UserPoolId = &a.AWSConfig.UserPoolId
// 		req, _ := tawsCIPPtr.AdminUpdateUserAttributesRequest(&tAdminUpdateUserAttributesInput)
// 		errorInfo.Error = req.Send()
// 	}
//
// 	return
// }

// ResetUserPassword - trigger one-time code to be set to user email.
// func (a *AWSHelper) ResetUserPassword(
// 	userName string,
// 	test bool,
// ) (errorInfo pi.ErrorInfo) {
//
// 	var (
// 		tAdminResetUserPasswordInput cognito.AdminResetUserPasswordInput
// 		tawsCIPPtr  *cognito.awsCIP
// 		tFunction, _, _, _           = runtime.Caller(0)
// 		tFunctionName                = runtime.FuncForPC(tFunction).Name()
// 		req                          *request.Request
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if userName == ctv.VAL_EMPTY {
// 		errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! AWS User Name: '%v'", userName))
// 		log.Println(errorInfo.Error)
// 	} else {
// 		tawsCIPPtr = cognito.New(a.SessionPtr)
// 		tAdminResetUserPasswordInput.Username = &userName
// 		tAdminResetUserPasswordInput.UserPoolId = &a.AWSConfig.UserPoolId
// 		if test == false {
// 			req, _ = tawsCIPPtr.AdminResetUserPasswordRequest(&tAdminResetUserPasswordInput)
// 			errorInfo.Error = req.Send()
// 		}
// 	}
//
// 	return
// }

// Private functions below here

// areAWSClaimsValid - Checks if email is verified and token is either an Id or Access token.
// func areAWSClaimsValid(
// 	FirestoreClientPtr *firestore.Client,
// 	subject, email, username, tokenUse string,
// 	emailVerified bool,
// ) bool {
//
// 	var (
// 		errorInfo          pi.ErrorInfo
// 		tDocumentPtr       *firestore.DocumentSnapshot
// 		tEmailInterface    interface{}
// 		tSubjectInterface  interface{}
// 		tUsernameInterface interface{}
// 	)
//
// 	if _, tDocumentPtr, errorInfo = coreFirestore.FindDocument(
// 		FirestoreClientPtr, ctv.DATASTORE_USERS, coreFirestore.NameValueQuery{
// 			FieldName:  ctv.FN_REQUESTOR_ID,
// 			FieldValue: subject,
// 		},
// 	); errorInfo.Error == nil {
// 		switch strings.ToUpper(tokenUse) {
// 		case ctv.TOKEN_TYPE_ID:
// 			if tSubjectInterface, errorInfo.Error = tDocumentPtr.DataAt(ctv.FN_REQUESTOR_ID); errorInfo.Error == nil {
// 				if tUsernameInterface, errorInfo.Error = tDocumentPtr.DataAt(ctv.FN_USERNAME); errorInfo.Error == nil {
// 					if tEmailInterface, errorInfo.Error = tDocumentPtr.DataAt(ctv.FN_EMAIL); errorInfo.Error == nil {
// 						if emailVerified && tSubjectInterface.(string) == subject && tEmailInterface.(string) == email && tUsernameInterface.(string) == username {
// 							return true
// 						}
// 					}
// 				}
// 			}
// 		case ctv.TOKEN_TYPE_ACCESS:
// 			if tSubjectInterface, errorInfo.Error = tDocumentPtr.DataAt(ctv.FN_REQUESTOR_ID); errorInfo.Error == nil {
// 				if tUsernameInterface, errorInfo.Error = tDocumentPtr.DataAt(ctv.FN_USERNAME); errorInfo.Error == nil {
// 					if emailVerified && tSubjectInterface.(string) == subject && tUsernameInterface.(string) == username {
// 						return true
// 					}
// 				}
// 			}
// 		}
// 	}
//
// 	return false
// }

// convertKey - decodes, processes, and returns the public key.
// NOTE: does not follow the errorInfo format because it is called by a function
// that only allows error to be returned.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func convertKey(rawE, rawN string) (
	publicKey *rsa.PublicKey,
	err error,
) {

	var (
		decodedN      []byte
		decodedBase64 []byte
		ndata         []byte
	)

	if decodedBase64, err = base64.RawURLEncoding.DecodeString(rawE); err == nil {
		if len(decodedBase64) < 4 {
			ndata = make([]byte, 4)
			copy(ndata[4-len(decodedBase64):], decodedBase64)
			decodedBase64 = ndata
		}
		publicKey = &rsa.PublicKey{
			N: &big.Int{},
			E: int(binary.BigEndian.Uint32(decodedBase64[:])),
		}
		if decodedN, err = base64.RawURLEncoding.DecodeString(rawN); err == nil {
			publicKey.N.SetBytes(decodedN)
		}
	}

	return
}

// getPublicKeySet - gets the public keys for AWS JWTs
//
//	Customer Messages: None
//	Errors: ErrRequiredArgumentMissing, ErrHTTPRequestFailed, http.Get or io.ReadAll or json.Unmarshal returned error
//	Verifications: None
func getPublicKeySet(keySetURL string) (
	keySet keySet,
	errorInfo pi.ErrorInfo,
) {

	var (
		tJWKS              map[string]interface{}
		tKey               key
		tKeySetResponsePtr *http.Response
		tKeyData           []byte
	)

	if keySetURL == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_MISSING_PARAMETER, ctv.FN_URL))
		return
	}

	if tKeySetResponsePtr, errorInfo.Error = http.Get(keySetURL); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, ctv.VAL_EMPTY)
		return
	}
	defer tKeySetResponsePtr.Body.Close()

	if tKeySetResponsePtr.StatusCode != ctv.HTTP_STATUS_200 {
		errorInfo = pi.NewErrorInfo(
			pi.ErrHTTPRequestFailed,
			fmt.Sprintf("%v%v - %v%v", ctv.TXT_HTTP_STATUS, tKeySetResponsePtr.StatusCode, ctv.FN_KEYSET_URL, keySetURL),
		)
		return
	}

	if errorInfo.Error = json.NewDecoder(tKeySetResponsePtr.Body).Decode(&tJWKS); errorInfo.Error != nil {
		return
	}
	if keys, ok := tJWKS["keys"].([]interface{}); ok {
		for _, key := range keys {
			if tKeyData, errorInfo.Error = json.Marshal(key); errorInfo.Error != nil {
				errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_PUBLIC_KEY, key))
				return
			}
			if errorInfo.Error = json.Unmarshal(tKeyData, &tKey); errorInfo.Error != nil {
				errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_PUBLIC_KEY, tKeyData))
				return
			}
			keySet.Keys = append(keySet.Keys, tKey) // Assuming "kid" is the key ID and "n" is the public key value
		}
		return
	}

	errorInfo = pi.NewErrorInfo(pi.ErrExtractKeysFailure, fmt.Sprintf("%v%v", ctv.TXT_PUBLIC_KEY, keySetURL))

	return
}

// func getTokenClaims(
// 	a *AWSHelper,
// 	tokenType, token string,
// ) (
// 	claimsPtr *Claims,
// 	errorInfo pi.ErrorInfo,
// ) {
//
// 	return a.ParseAWSJWTWithClaims(tokenType, token)
// }

// func isTokenValid(
// 	firestoreClientPtr *firestore.Client,
// 	a *AWSHelper,
// 	tokenType, token string,
// ) bool {
//
// 	var (
// 		errorInfo  pi.ErrorInfo
// 		tClaimsPtr *Claims
// 	)
//
// 	a.tokenType = tokenType
// 	if tClaimsPtr, errorInfo = getTokenClaims(a, tokenType, token); errorInfo.Error == nil {
// 		switch strings.ToUpper(tClaimsPtr.TokenUse) {
// 		case ctv.TOKEN_TYPE_ID:
// 			return areAWSClaimsValid(
// 				firestoreClientPtr,
// 				tClaimsPtr.Subject,
// 				tClaimsPtr.Email,
// 				tClaimsPtr.CognitoUsername,
// 				tClaimsPtr.TokenUse,
// 				tClaimsPtr.EmailVerified,
// 			)
// 		case ctv.TOKEN_TYPE_ACCESS:
// 			return areAWSClaimsValid(firestoreClientPtr, tClaimsPtr.Subject, ctv.VAL_EMPTY, tClaimsPtr.UserName, tClaimsPtr.TokenUse, true)
// 		}
// 	}
//
// 	return false
// }
