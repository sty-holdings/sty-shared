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
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCIP "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	awsCT "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// NewAWSConfig - reads the SDKs default external configurations, and populates an AWS Config with the values from the external configurations.
//
//	Customer Messages: None
//	Errors: ErrEnvironmentInvalid, anything awsConfig.LoadDefaultConfig or getPublicKeySet returns,
//	Verifications: None
func NewAWSConfig(environment string) (
	session AWSSession,
	errorInfo pi.ErrorInfo,
) {

	switch strings.ToLower(strings.Trim(environment, ctv.SPACES_ONE)) {
	case ctv.ENVIRONMENT_PRODUCTION:
		session.STYConfig = styConfig
	case ctv.ENVIRONMENT_DEVELOPMENT:
		session.STYConfig = styConfigDevelopment
	case ctv.ENVIRONMENT_LOCAL:
		session.STYConfig = styConfigLocal
	default:
		errorInfo = pi.NewErrorInfo(pi.ErrEnvironmentInvalid, fmt.Sprintf("%v%v", ctv.TXT_EVIRONMENT, environment))
	}

	if session.BaseConfig, errorInfo.Error = awsConfig.LoadDefaultConfig(awsCTX, awsConfig.WithRegion(session.STYConfig.Region)); errorInfo.
		Error != nil {
		errorInfo = pi.NewErrorInfo(pi.ErrServiceFailedAWS, "Failed to create an AWS Session.")
		return
	}

	session.KeySetURL = fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", session.STYConfig.Region, session.STYConfig.UserPoolId,
	)
	session.KeySet, errorInfo = getPublicKeySet(session.KeySetURL)

	return
}

// GetParameters - will return System Manager parameters
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func GetParameters(
	loginType, username string,
	password *string,
	session AWSSession,
) (
	tokens map[string]string,
	errorInfo pi.ErrorInfo,
) {

	if loginType == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_LOGIN_TYPE, loginType))
		return
	}
	if username == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_USERNAME, username))
		return
	}
	if loginType == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_PASSWORD, ctv.TXT_PROTECTED))
		return
	}

	csrp, _ := NewCognitoLogin(username, session.STYConfig.UserPoolId, session.STYConfig.ClientId, password, nil)

	svc := awsCIP.NewFromConfig(session.BaseConfig)

	// initiate auth
	resp, err := svc.InitiateAuth(
		context.Background(), &awsCIP.InitiateAuthInput{
			AuthFlow:       awsCT.AuthFlowType(loginType),
			ClientId:       aws.String(csrp.GetClientId()),
			AuthParameters: csrp.GetAuthParams(),
		},
	)
	if err != nil {
		panic(err)
	}

	tokens = make(map[string]string) // This is used for either awsCT.AuthFlowTypeUserPasswordAuth or awsCT.AuthFlowTypeUserSrpAuth
	if loginType == string(awsCT.AuthFlowTypeUserPasswordAuth) {
		tokens["access"] = *resp.AuthenticationResult.AccessToken
		tokens["id"] = *resp.AuthenticationResult.IdToken
		tokens["refresh"] = *resp.AuthenticationResult.RefreshToken
		return
	}

	// respond to password verifier challenge
	if resp.ChallengeName == awsCT.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		x, err := svc.RespondToAuthChallenge(
			context.Background(), &awsCIP.RespondToAuthChallengeInput{
				ChallengeName:      awsCT.ChallengeNameTypePasswordVerifier,
				ChallengeResponses: challengeResponses,
				ClientId:           aws.String(csrp.GetClientId()),
			},
		)
		if err != nil {
			panic(err)
		}
		tokens["access"] = *x.AuthenticationResult.AccessToken
		tokens["id"] = *x.AuthenticationResult.IdToken
		tokens["refresh"] = *x.AuthenticationResult.RefreshToken
	}

	return
}

// Login - will validate and return tokens, if the login is successful.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func Login(
	loginType, username string,
	password *string,
	session AWSSession,
) (
	tokens map[string]string,
	errorInfo pi.ErrorInfo,
) {

	if loginType == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_LOGIN_TYPE, loginType))
		return
	}
	if username == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_USERNAME, username))
		return
	}
	if loginType == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_PASSWORD, ctv.TXT_PROTECTED))
		return
	}

	cognitoLogin, _ := NewCognitoLogin(username, session.STYConfig.UserPoolId, session.STYConfig.ClientId, password, nil)

	service := awsCIP.NewFromConfig(session.BaseConfig)

	// initiate auth
	resp, err := service.InitiateAuth(
		context.Background(), &awsCIP.InitiateAuthInput{
			AuthFlow:       awsCT.AuthFlowType(loginType),
			ClientId:       aws.String(cognitoLogin.GetClientId()),
			AuthParameters: cognitoLogin.GetAuthParams(),
		},
	)
	if err != nil {
		panic(err)
	}

	tokens = make(map[string]string) // This is used for either awsCT.AuthFlowTypeUserPasswordAuth or awsCT.AuthFlowTypeUserSrpAuth
	if loginType == string(awsCT.AuthFlowTypeUserPasswordAuth) {
		tokens["access"] = *resp.AuthenticationResult.AccessToken
		tokens["id"] = *resp.AuthenticationResult.IdToken
		tokens["refresh"] = *resp.AuthenticationResult.RefreshToken
		return
	}

	// respond to password verifier challenge
	if resp.ChallengeName == awsCT.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := cognitoLogin.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		x, err := service.RespondToAuthChallenge(
			context.Background(), &awsCIP.RespondToAuthChallengeInput{
				ChallengeName:      awsCT.ChallengeNameTypePasswordVerifier,
				ChallengeResponses: challengeResponses,
				ClientId:           aws.String(cognitoLogin.GetClientId()),
			},
		)
		if err != nil {
			panic(err)
		}
		tokens["access"] = *x.AuthenticationResult.AccessToken
		tokens["id"] = *x.AuthenticationResult.IdToken
		tokens["refresh"] = *x.AuthenticationResult.RefreshToken
	}

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
// 		tAdminConfirmSignUpInput.UserPoolId = &a.AWSConfig.UserPoolId
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

// ParseAWSJWT - will return the claims, if any, or an err if the AWS JWT is invalid.
// This will parse ID and Access tokens. Refresh token are not support and nothing is returned.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func ParseAWSJWT(
	session AWSSession,
	tokenType, token string,
) (
	claims jwt.Claims,
	tokenValuePtr *jwt.Token,
	errorInfo pi.ErrorInfo,
) {

	if len(session.KeySet.Keys) == ctv.VAL_ZERO {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, ctv.TXT_KEY_SET_MISSING)
		return
	}
	if token == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_TOKEN, ctv.TXT_PROTECTED))
		return
	}

	if tokenType == ctv.TOKEN_TYPE_REFRESH {
		return
	}

	for i := 0; i < len(session.KeySet.Keys); i++ {
		if tokenValuePtr, errorInfo.Error = jwt.ParseWithClaims(
			token, jwt.MapClaims{}, func(token *jwt.Token) (
				key interface{},
				err error,
			) {
				key, err = convertKey(session.KeySet.Keys[i].E, session.KeySet.Keys[i].N) // ID
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

	errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_TOKEN, ctv.TXT_PROTECTED))

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
	keySet KeySet,
	errorInfo pi.ErrorInfo,
) {

	var (
		tJWKS              map[string]interface{}
		tKey               Key
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
			fmt.Sprintf("%v%v - %v%v", ctv.TXT_HTTP_STATUS, tKeySetResponsePtr.StatusCode, ctv.FN_URL, keySetURL),
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
