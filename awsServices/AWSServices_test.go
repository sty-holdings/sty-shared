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
	"testing"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

func TestAssumeRole(tPtr *testing.T) {

	type arguments struct {
		loginType          string
		password           *string
		shouldBeAuthorized bool
		username           string
	}

	var (
		environment = ctv.ENVIRONMENT_PRODUCTION
		errorInfo   pi.ErrorInfo
		gotError    bool
		password    = "Yidiao09#1"
		sessionPtr  *AWSSession
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with username/password",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_PASSWORD_AUTH,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with SRP",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_SRP,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				sessionPtr, errorInfo = NewAWSConfig(environment)
				_, errorInfo = Login(ts.arguments.loginType, ts.arguments.username, ts.arguments.password, sessionPtr)
				if _, errorInfo = AssumeRole(sessionPtr, ctv.VAL_EMPTY); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Error(errorInfo.Error.Error())
				}
			},
		)
	}
}

func TestGetIdentityCredentials(tPtr *testing.T) {

	type arguments struct {
		loginType          string
		password           *string
		shouldBeAuthorized bool
		username           string
	}

	var (
		environment = ctv.ENVIRONMENT_PRODUCTION
		errorInfo   pi.ErrorInfo
		gotError    bool
		password    = "Yidiao09#1"
		sessionPtr  *AWSSession
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with username/password",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_PASSWORD_AUTH,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with SRP",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_SRP,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				sessionPtr, errorInfo = NewAWSConfig(environment)
				_, errorInfo = Login(ts.arguments.loginType, ts.arguments.username, ts.arguments.password, sessionPtr)
				_, errorInfo = GetId(sessionPtr, ctv.VAL_EMPTY, ctv.VAL_EMPTY)
				if _, errorInfo = GetIdentityCredentials(sessionPtr, ctv.VAL_EMPTY); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Error(errorInfo.Error.Error())
				}
			},
		)
	}
}

func TestGetId(tPtr *testing.T) {

	type arguments struct {
		loginType          string
		password           *string
		shouldBeAuthorized bool
		username           string
	}

	var (
		environment = ctv.ENVIRONMENT_PRODUCTION
		errorInfo   pi.ErrorInfo
		gotError    bool
		password    = "Yidiao09#1"
		sessionPtr  *AWSSession
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with username/password",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_PASSWORD_AUTH,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with SRP",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_SRP,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				sessionPtr, errorInfo = NewAWSConfig(environment)
				_, errorInfo = Login(ts.arguments.loginType, ts.arguments.username, ts.arguments.password, sessionPtr)
				if _, errorInfo = GetId(sessionPtr, ctv.VAL_EMPTY, ctv.VAL_EMPTY); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Error(errorInfo.Error.Error())
				}
			},
		)
	}
}

func TestGetParameters(tPtr *testing.T) {

	type arguments struct {
		loginType          string
		password           *string
		shouldBeAuthorized bool
		username           string
	}

	var (
		environment = ctv.ENVIRONMENT_PRODUCTION
		errorInfo   pi.ErrorInfo
		gotError    bool
		password    = "Yidiao09#1"
		sessionPtr  *AWSSession
		tokens      map[string]string
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with username/password",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_PASSWORD_AUTH,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with SRP",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_SRP,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				sessionPtr, errorInfo = NewAWSConfig(environment)
				tokens, errorInfo = Login(ts.arguments.loginType, ts.arguments.username, ts.arguments.password, sessionPtr)
				if _, errorInfo = GetParameters(sessionPtr, tokens[ctv.TOKEN_TYPE_ID], "ai2-development-nats-token"); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Error(errorInfo.Error.Error())
				}
			},
		)
	}
}

func TestLogin(tPtr *testing.T) {

	type arguments struct {
		loginType          string
		password           *string
		shouldBeAuthorized bool
		username           string
	}

	var (
		environment = ctv.ENVIRONMENT_PRODUCTION
		errorInfo   pi.ErrorInfo
		gotError    bool
		password    = "Yidiao09#1"
		sessionPtr  *AWSSession
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with username/password",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_PASSWORD_AUTH,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with SRP",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_SRP,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				sessionPtr, errorInfo = NewAWSConfig(environment)
				if _, errorInfo = Login(ts.arguments.loginType, ts.arguments.username, ts.arguments.password, sessionPtr); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Error(errorInfo.Error.Error())
				}
			},
		)
	}
}

// Part of run_AWS_No_Token_Test list
// func TestAWSHelper_ConfirmUser(tPtr *testing.T) {
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		tAWSHelper        AWSHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 	)
//
// 	tAWSHelper, _ = NewAWSSession(ctv.TEST_AWS_INFORMATION_FQN)
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			if errorInfo = tAWSHelper.ConfirmUser(ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE); errorInfo.Error != nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 			}
// 			if errorInfo = tAWSHelper.ConfirmUser(ctv.EMPTY); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 		},
// 	)
// }

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
// func TestAWSHelper_GetRequestorEmailPhoneFromIdTokenClaims(tPtr *testing.T) {
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myAWS             AWSHelper
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 	)
//
// 	myAWS, myFireBase = StartTest()
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			//  Positive Test - Successful
// 			if _, _, _, errorInfo = myAWS.GetRequestorEmailPhoneFromIdTokenClaims(
// 				myFireBase.FirestoreClientPtr,
// 				string(testingIdTokenValid),
// 			); errorInfo.Error != nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 			}
// 			if _, _, _, errorInfo = myAWS.GetRequestorEmailPhoneFromIdTokenClaims(
// 				myFireBase.FirestoreClientPtr,
// 				ctv.TEST_TOKEN_INVALID,
// 			); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 			if _, _, _, errorInfo = myAWS.GetRequestorEmailPhoneFromIdTokenClaims(myFireBase.FirestoreClientPtr, ctv.EMPTY); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 		},
// 	)
//
// 	StopTest(myFireBase)
// }

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
// func TestAWSHelper_GetRequestorFromAccessTokenClaims(tPtr *testing.T) {
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myAWS             AWSHelper
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 	)
//
// 	myAWS, myFireBase = StartTest()
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			//  Positive Test - Successful
// 			if _, errorInfo = myAWS.GetRequestorFromAccessTokenClaims(
// 				myFireBase.FirestoreClientPtr,
// 				string(testingAccessTokenValid),
// 			); errorInfo.Error != nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 			}
// 			if _, errorInfo = myAWS.GetRequestorFromAccessTokenClaims(myFireBase.FirestoreClientPtr, ctv.TEST_TOKEN_INVALID); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 			if _, errorInfo = myAWS.GetRequestorFromAccessTokenClaims(myFireBase.FirestoreClientPtr, ctv.EMPTY); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 		},
// 	)
//
// 	StopTest(myFireBase)
// }

func TestNewAWSSession(tPtr *testing.T) {

	type arguments struct {
		config      AWSConfig
		environment string
	}

	var (
		errorInfo pi.ErrorInfo
		gotError  bool
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "local config",
			arguments: arguments{
				config:      styConfigLocal,
				environment: ctv.ENVIRONMENT_LOCAL,
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "development config",
			arguments: arguments{
				config:      styConfigDevelopment,
				environment: ctv.ENVIRONMENT_DEVELOPMENT,
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "production config",
			arguments: arguments{
				config:      styConfig,
				environment: ctv.ENVIRONMENT_PRODUCTION,
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				if _, errorInfo = NewAWSConfig(ts.arguments.environment); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Error(errorInfo.Error.Error())
				}
			},
		)
	}
}

// Requires a JWT. You can get a token two ways:
// 1) You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View
// Hosted UI to login. This will output an access and id token for the user.
// 2) Call the AWSServices Login function before each test needing a token
func TestParseAWSJWT(tPtr *testing.T) {

	type arguments struct {
		loginType          string
		password           *string
		shouldBeAuthorized bool
		username           string
	}

	var (
		environment = ctv.ENVIRONMENT_PRODUCTION
		errorInfo   pi.ErrorInfo
		gotError    bool
		password    = "Yidiao09#1"
		sessionPtr  *AWSSession
		tokens      = make(map[string]string)
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with username/password",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_PASSWORD_AUTH,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "login with SRP",
			arguments: arguments{
				loginType:          ctv.AUTH_USER_SRP,
				password:           &password,
				shouldBeAuthorized: true,
				username:           "scott@yackofamily.com",
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				sessionPtr, errorInfo = NewAWSConfig(environment)
				tokens, errorInfo = Login(ts.arguments.loginType, ts.arguments.username, ts.arguments.password, sessionPtr)
				for tokenType, token := range tokens {
					if ParseAWSJWT(sessionPtr, tokenType, token); errorInfo.Error != nil {
						gotError = true
					} else {
						gotError = false
					}
					if gotError != ts.wantError {
						tPtr.Error(errorInfo.Error.Error())
					}
				}
			},
		)
	}
}

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
// func TestAWSHelper_ParseJWT(tPtr *testing.T) {
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myAWS             AWSHelper
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 	)
//
// 	myAWS, myFireBase = StartTest()
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			if _, errorInfo = myAWS.ParseAWSJWT(string(testingAccessTokenValid)); errorInfo.Error != nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 			}
// 			if _, errorInfo = myAWS.ParseAWSJWT(ctv.EMPTY); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 		},
// 	)
//
// 	StopTest(myFireBase)
// }

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
// func TestAWSHelper_PullCognitoUserInfo(tPtr *testing.T) {
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myAWS             AWSHelper
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 	)
//
// 	myAWS, myFireBase = StartTest()
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			if _, errorInfo = myAWS.PullCognitoUserInfo(ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE); errorInfo.Error != nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 			}
// 			if _, errorInfo = myAWS.PullCognitoUserInfo(ctv.EMPTY); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 		},
// 	)
//
// 	StopTest(myFireBase)
// }

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
//
//	The actual reset will be bypassed because the resetByPass is set to true
// func TestAWSHelper_ResetUserPassword(tPtr *testing.T) {
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myAWS             AWSHelper
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		resetByPass       = true
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 	)
//
// 	myAWS, myFireBase = StartTest()
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			if errorInfo = myAWS.ResetUserPassword(ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE, resetByPass); errorInfo.Error != nil {
// 				if errorInfo.Error.Error() == pi.ATTEMPTS_EXCEEDED {
// 					fmt.Println(pi.ATTEMPTS_EXCEEDED)
// 				} else {
// 					tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 				}
// 			}
// 			if errorInfo = myAWS.ResetUserPassword(ctv.EMPTY, resetByPass); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 		},
// 	)
//
// 	StopTest(myFireBase)
// }

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
// func TestAWSHelper_UpdateAWSEmailVerifyFlag(tPtr *testing.T) {
// 	//
// 	// NOTE: The Id and Access token must match the username in ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE
// 	//
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myAWS             AWSHelper
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 		tUsername         = ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE
// 	)
//
// 	myAWS, myFireBase = StartTest()
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			if errorInfo = myAWS.UpdateAWSEmailVerifyFlag(tUsername); errorInfo.Error != nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 			}
// 			if errorInfo = myAWS.UpdateAWSEmailVerifyFlag(ctv.EMPTY); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 		},
// 	)
//
// 	StopTest(myFireBase)
// }

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
// func TestAWSHelper_ValidAWSJWT(tPtr *testing.T) {
//
// 	type arguments struct {
// 		tokenType string
// 		token     string
// 	}
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myAWS             AWSHelper
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 		tToken            string
// 		tValid            bool
// 	)
//
// 	tests := []struct {
// 		name          string
// 		arguments     arguments
// 		shouldBeValid bool
// 	}{
// 		{
// 			name: "Positive Case: Access Successful!",
// 			arguments: arguments{
// 				tokenType: ctv.TOKEN_TYPE_ACCESS,
// 				token:     "valid",
// 			},
// 			shouldBeValid: true,
// 		},
// 		{
// 			name: "Positive Case: Id Successful!",
// 			arguments: arguments{
// 				tokenType: ctv.TOKEN_TYPE_ID,
// 				token:     "valid",
// 			},
// 			shouldBeValid: true,
// 		},
// 		{
// 			name: "Negative Case: Access invalid!",
// 			arguments: arguments{
// 				tokenType: ctv.TOKEN_TYPE_ACCESS,
// 				token:     "invalid",
// 			},
// 			shouldBeValid: false,
// 		},
// 		{
// 			name: "Negative Case: Id invalid!",
// 			arguments: arguments{
// 				tokenType: ctv.TOKEN_TYPE_ID,
// 				token:     "invalid",
// 			},
// 			shouldBeValid: false,
// 		},
// 		{
// 			name: "Negative Case: Access missing!",
// 			arguments: arguments{
// 				tokenType: ctv.TOKEN_TYPE_ACCESS,
// 				token:     "missing",
// 			},
// 			shouldBeValid: false,
// 		},
// 		{
// 			name: "Negative Case: Id missing!",
// 			arguments: arguments{
// 				tokenType: ctv.TOKEN_TYPE_ID,
// 				token:     "missing",
// 			},
// 			shouldBeValid: false,
// 		},
// 	}
//
// 	myAWS, myFireBase = StartTest()
//
// 	for _, ts := range tests {
// 		tPtr.Run(
// 			ts.name, func(t *testing.T) {
// 				tToken = getToken(ts.arguments.tokenType, ts.arguments.token)
// 				if tValid, errorInfo = myAWS.ValidAWSJWT(myFireBase.FirestoreClientPtr, ts.arguments.tokenType, tToken); tValid != ts.shouldBeValid {
// 					tPtr.Error(tFunctionName, ts.name, errorInfo, fmt.Sprintf("Expected the token to be %v and it was %v", ts.shouldBeValid, tValid))
// 				}
// 			},
// 		)
// 	}
//
// 	StopTest(myFireBase)
// }

// Part of run_AWS_No_Token_Test list
// func TestGetPublicKeySet(tPtr *testing.T) {
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 		tKeySetURL        = fmt.Sprintf(ctv.TEST_AWS_KEYSET_URL, ctv.TEST_USER_POOL_ID)
// 	)
//
// 	tPtr.Run(
// 		tFunctionName, func(tPtr *testing.T) {
// 			if _, errorInfo = getPublicKeySet(tKeySetURL); errorInfo.Error != nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, "nil", errorInfo.Error.Error())
// 			}
// 			if _, errorInfo = getPublicKeySet(ctv.EMPTY); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, errorInfo.Error.Error(), "nil")
// 			}
// 			if _, errorInfo = getPublicKeySet(ctv.TEST_URL_INVALID); errorInfo.Error == nil {
// 				tPtr.Errorf("%v Failed: Was expecting an err of %v but got %v.", tFunctionName, `errorInfo.Error.Error()`, "nil")
// 			}
// 		},
// 	)
// }

// Requires updated access token. You can use Cognito > User pools > App integration > App clients and analytics > {app name} > Hosted UI > View Hosted UI
// to login. This will output an access and id token for the user.
// func TestValidAWSClaims(tPtr *testing.T) {
//
// 	type arguments struct {
// 		subject       string
// 		email         string
// 		username      string
// 		emailVerified bool // emailVerified is only checked for ctv.TOKEN_TYPE_ID
// 		tokenUse      string
// 	}
//
// 	var (
// 		errorInfo         pi.ErrorInfo
// 		myFireBase        coreHelpers.FirebaseFirestoreHelper
// 		function, _, _, _ = runtime.Caller(0)
// 		tFunctionName     = runtime.FuncForPC(function).Name()
// 		tValid            bool
// 	)
//
// 	tests := []struct {
// 		name          string
// 		arguments     arguments
// 		shouldBeValid bool
// 	}{
// 		{
// 			name: "Positive Case: Successful Id Token!",
// 			arguments: arguments{
// 				subject:       ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID,
// 				email:         ctv.TEST_USER_EMAIL,
// 				username:      ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE,
// 				emailVerified: true,
// 				tokenUse:      ctv.TOKEN_TYPE_ID,
// 			},
// 			shouldBeValid: true,
// 		},
// 		{
// 			name: "Positive Case: Successful Access Token!",
// 			arguments: arguments{
// 				subject:       ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID,
// 				email:         ctv.TEST_USER_EMAIL,
// 				username:      ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE,
// 				emailVerified: true,
// 				tokenUse:      ctv.TOKEN_TYPE_ACCESS,
// 			},
// 			shouldBeValid: true,
// 		},
// 		{
// 			name: "Negative Case: Email not verified!",
// 			arguments: arguments{
// 				subject:       ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID,
// 				email:         ctv.TEST_USER_EMAIL,
// 				username:      ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE,
// 				emailVerified: false,
// 				tokenUse:      ctv.TOKEN_TYPE_ID,
// 			},
// 			shouldBeValid: false,
// 		},
// 		{
// 			name: "Negative Case: Token type missing!",
// 			arguments: arguments{
// 				subject:       ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID,
// 				email:         ctv.TEST_USER_EMAIL,
// 				username:      ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE,
// 				emailVerified: true,
// 				tokenUse:      ctv.EMPTY,
// 			},
// 			shouldBeValid: false,
// 		},
// 		{
// 			name: "Negative Case: subject is missing!",
// 			arguments: arguments{
// 				subject:       ctv.EMPTY,
// 				email:         ctv.TEST_USER_EMAIL,
// 				username:      ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE,
// 				emailVerified: true,
// 				tokenUse:      ctv.EMPTY,
// 			},
// 			shouldBeValid: false,
// 		},
// 		{
// 			name: "Negative Case: email is missing!",
// 			arguments: arguments{
// 				subject:       ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID,
// 				email:         ctv.EMPTY,
// 				username:      ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE,
// 				emailVerified: true,
// 				tokenUse:      ctv.EMPTY,
// 			},
// 			shouldBeValid: false,
// 		},
// 		{
// 			name: "Negative Case: username is missing!",
// 			arguments: arguments{
// 				subject:       ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID,
// 				email:         ctv.TEST_USER_EMAIL,
// 				username:      ctv.EMPTY,
// 				emailVerified: true,
// 				tokenUse:      ctv.EMPTY,
// 			},
// 			shouldBeValid: false,
// 		},
// 	}
//
// 	_, myFireBase = StartTest()
//
// 	for _, ts := range tests {
// 		tPtr.Run(
// 			ts.name, func(t *testing.T) {
// 				if tValid = areAWSClaimsValid(
// 					myFireBase.FirestoreClientPtr,
// 					ts.arguments.subject,
// 					ts.arguments.email,
// 					ts.arguments.username,
// 					ts.arguments.tokenUse,
// 					ts.arguments.emailVerified,
// 				); tValid != ts.shouldBeValid {
// 					tPtr.Error(tFunctionName, ts.name, errorInfo, fmt.Sprintf("Expected the token to be %v and it was %v", ts.shouldBeValid, tValid))
// 				}
// 			},
// 		)
// 	}
//
// 	StopTest(myFireBase)
// }
