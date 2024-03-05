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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// NewAWSSession
//
//	Customer Messages: None
//	Errors:
//	Verifications: None
func NewAWSSession() (
	awsHelper AWSHelper,
	errorInfo pi.ErrorInfo,
) {

	if awsHelper.SessionPtr, errorInfo.Error = awsSession.NewSessionWithOptions(
		awsSession.Options{
			Config: aws.Config{
				Region: aws.String(awsConfig.Region),
			},
		},
	); errorInfo.Error != nil {
		return
	}

	awsHelper.KeySetURL = fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsConfig.Region, awsConfig.UserPoolId,
	)
	awsHelper.KeySet, errorInfo = getPublicKeySet(awsHelper.KeySetURL)

	return
}

// // ConfirmUser - mark the AWS user as confirmed
// func (a *AWSHelper) ConfirmUser(userName string) (errorInfo pi.ErrorInfo) {
//
// 	var (
// 		tAdminConfirmSignUpInput    cognito.AdminConfirmSignUpInput
// 		tCognitoIdentityProviderPtr *cognito.CognitoIdentityProvider
// 		tFunction, _, _, _          = runtime.Caller(0)
// 		tFunctionName               = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	if userName == ctv.VAL_EMPTY {
// 		errorInfo.Error = pi.ErrRequiredArgumentMissing
// 		log.Println(errorInfo.Error)
// 	} else {
// 		tCognitoIdentityProviderPtr = cognito.New(a.SessionPtr)
// 		tAdminConfirmSignUpInput.Username = &userName
// 		tAdminConfirmSignUpInput.UserPoolId = &a.AWSConfig.UserPoolId
// 		if _, errorInfo.Error = tCognitoIdentityProviderPtr.AdminConfirmSignUp(&tAdminConfirmSignUpInput); errorInfo.Error != nil {
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
// 		if token == ctv.EMPTY {
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

// ParseAWSJWTWithClaims - will return an err if the AWS JWT is invalid.
// func (a *AWSHelper) ParseAWSJWTWithClaims(tokenType, tokenString string) (
// 	claimsPtr *Claims,
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
// 	if tokenString == ctv.EMPTY {
// 		errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! Token: '%v'", tokenString))
// 		log.Println(errorInfo.Error)
// 	} else {
// 		if _, errorInfo.Error = jwt.ParseWithClaims(
// 			tokenString,
// 			&Claims{},
// 			func(token *jwt.Token) (
// 				key interface{},
// 				err error,
// 			) {
// 				switch strings.ToUpper(tokenType) {
// 				case ctv.TOKEN_TYPE_ID:
// 					key, err = convertKey(a.KeySet.Keys[0].E, a.KeySet.Keys[0].N)
// 				case ctv.TOKEN_TYPE_ACCESS:
// 					key, err = convertKey(a.KeySet.Keys[1].E, a.KeySet.Keys[1].N)
// 				}
// 				claimsPtr = token.Claims.(*Claims)
// 				return
// 			},
// 		); errorInfo.Error != nil {
// 			log.Println(errorInfo.Error)
// 		}
// 	}
//
// 	return
// }

// ParseAWSJWT - will return an err if the AWS JWT is invalid.
// func (a *AWSHelper) ParseAWSJWT(tokenString string) (
// 	tTokenPtr *jwt.Token,
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
// 	if tokenString == ctv.EMPTY {
// 		errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! Token: '%v'", tokenString))
// 		fmt.Println(errorInfo.Error)
// 	} else {
// 		tTokenPtr, errorInfo.Error = jwt.Parse(
// 			tokenString,
// 			func(token *jwt.Token) (
// 				key interface{},
// 				err error,
// 			) {
// 				key, err = convertKey(a.KeySet.Keys[1].E, a.KeySet.Keys[1].N)
// 				return
// 			},
// 		)
// 	}
//
// 	return
// }

// PullCognitoUserInfo - pull user information from the Cognito user pool.
// func (a *AWSHelper) PullCognitoUserInfo(username string) (
// 	userData map[string]interface{},
// 	errorInfo pi.ErrorInfo,
// ) {
//
// 	var (
// 		tAdminGetUserInput          cognito.AdminGetUserInput
// 		tAdminGetUserOutputPtr      *cognito.AdminGetUserOutput
// 		tCognitoIdentityProviderPtr *cognito.CognitoIdentityProvider
// 		tFunction, _, _, _          = runtime.Caller(0)
// 		tFunctionName               = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if username == ctv.EMPTY {
// 		errorInfo.Error = pi.ErrRequiredArgumentMissing
// 		pi.PrintError(errorInfo)
// 	} else {
// 		tCognitoIdentityProviderPtr = cognito.New(a.SessionPtr)
// 		if tCognitoIdentityProviderPtr == nil {
// 			errorInfo.FileName, errorInfo.ErrorLineNumber = pi.GetFileLineNumber()
// 			errorInfo.Error = pi.ErrServiceFailedAWS
// 			pi.PrintError(errorInfo)
// 		} else {
// 			// Set up the request to get the user
// 			tAdminGetUserInput.UserPoolId = &a.AWSConfig.UserPoolId
// 			tAdminGetUserInput.Username = &username
// 			// Make the request to get the user
// 			if tAdminGetUserOutputPtr, errorInfo.Error = tCognitoIdentityProviderPtr.AdminGetUser(&tAdminGetUserInput); errorInfo.Error == nil {
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
// 	if token == ctv.EMPTY {
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
// 		tCognitoIdentityProviderPtr     *cognito.CognitoIdentityProvider
// 		tFunction, _, _, _              = runtime.Caller(0)
// 		tFunctionName                   = runtime.FuncForPC(tFunction).Name()
// 		tName                           string
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if username == ctv.EMPTY {
// 		errorInfo.Error = pi.ErrRequiredArgumentMissing
// 	} else {
// 		tCognitoIdentityProviderPtr = cognito.New(a.SessionPtr)
// 		tName = ctv.FN_EMAIL_VERIFIED // This is required because go doesn't support pointers to ctv.
// 		tAttributeType = cognito.AttributeType{
// 			Name:  &tName,
// 			Value: &tTrueString,
// 		}
// 		tAttributeTypePtrs = append(tAttributeTypePtrs, &tAttributeType)
// 		tAdminUpdateUserAttributesInput.UserAttributes = tAttributeTypePtrs
// 		tAdminUpdateUserAttributesInput.Username = &username
// 		tAdminUpdateUserAttributesInput.UserPoolId = &a.AWSConfig.UserPoolId
// 		req, _ := tCognitoIdentityProviderPtr.AdminUpdateUserAttributesRequest(&tAdminUpdateUserAttributesInput)
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
// 		tCognitoIdentityProviderPtr  *cognito.CognitoIdentityProvider
// 		tFunction, _, _, _           = runtime.Caller(0)
// 		tFunctionName                = runtime.FuncForPC(tFunction).Name()
// 		req                          *request.Request
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	if userName == ctv.EMPTY {
// 		errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! AWS User Name: '%v'", userName))
// 		log.Println(errorInfo.Error)
// 	} else {
// 		tCognitoIdentityProviderPtr = cognito.New(a.SessionPtr)
// 		tAdminResetUserPasswordInput.Username = &userName
// 		tAdminResetUserPasswordInput.UserPoolId = &a.AWSConfig.UserPoolId
// 		if test == false {
// 			req, _ = tCognitoIdentityProviderPtr.AdminResetUserPasswordRequest(&tAdminResetUserPasswordInput)
// 			errorInfo.Error = req.Send()
// 		}
// 	}
//
// 	return
// }

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

// convertKey - does not follow the errorInfo format because it is called by a function that only allows error to be returned
// func convertKey(rawE, rawN string) (
// 	publicKey *rsa.PublicKey,
// 	err error,
// ) {
//
// 	var (
// 		decodedN      []byte
// 		decodedBase64 []byte
// 		ndata         []byte
// 	)
//
// 	decodedBase64, err = base64.RawURLEncoding.DecodeString(rawE)
//
// 	if err == nil {
// 		if len(decodedBase64) < 4 {
// 			ndata = make([]byte, 4)
// 			copy(ndata[4-len(decodedBase64):], decodedBase64)
// 			decodedBase64 = ndata
// 		}
// 		publicKey = &rsa.PublicKey{
// 			N: &big.Int{},
// 			E: int(binary.BigEndian.Uint32(decodedBase64[:])),
// 		}
// 		if decodedN, err = base64.RawURLEncoding.DecodeString(rawN); err == nil {
// 			publicKey.N.SetBytes(decodedN)
// 		}
// 	}
//
// 	return
// }

// getPublicKeySet - gets the public key for AWS JWTs
//
//	Customer Messages: None
//	Errors: ErrRequiredArgumentMissing, ErrHTTPRequestFailed, http.Get or io.ReadAll or json.Unmarshal returned error
//	Verifications: None
func getPublicKeySet(keySetURL string) (
	keySet KeySet,
	errorInfo pi.ErrorInfo,
) {

	var (
		tBody      []byte
		tKeySetPtr *http.Response
	)

	if keySetURL == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_MISSING_PARAMETER, ctv.FN_URL))
		return
	}

	if tKeySetPtr, errorInfo.Error = http.Get(keySetURL); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, ctv.VAL_EMPTY)
		return
	}

	if tKeySetPtr.StatusCode != ctv.HTTP_STATUS_200 {
		errorInfo = pi.NewErrorInfo(
			pi.ErrHTTPRequestFailed,
			fmt.Sprintf("%v%v - %v%v", ctv.TXT_HTTP_STATUS, tKeySetPtr.StatusCode, ctv.FN_URL, keySetURL),
		)
		return
	}

	if tBody, errorInfo.Error = io.ReadAll(tKeySetPtr.Body); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_HTTP_BODY, tKeySetPtr.Body))
		return
	}

	if errorInfo.Error = json.Unmarshal(tBody, &keySet); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_HTTP_BODY, tKeySetPtr.Body))
		return
	}

	return
}

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
// 			return areAWSClaimsValid(firestoreClientPtr, tClaimsPtr.Subject, ctv.EMPTY, tClaimsPtr.UserName, tClaimsPtr.TokenUse, true)
// 		}
// 	}
//
// 	return false
// }

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
