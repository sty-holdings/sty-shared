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
	"fmt"
	"os"
	"runtime"
	"testing"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

var (
// Global variables here
)

func TestGenerateConfigFileSkeleton(tPtr *testing.T) {

	var (
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
	)

	tPtr.Run(
		tFunctionName, func(tPtr *testing.T) {
			GenerateConfigFileSkeleton(
				"NATS Connect Test",
				DEFAULT_SKELETON_CONFIG_FQD,
			)
		},
	)
}

func TestProcessBaseConfigFile(tPtr *testing.T) {

	type arguments struct {
		configFileName string
	}

	var (
		gotError           bool
		errorInfo          pi.ErrorInfo
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid config file",
			arguments: arguments{
				configFileName: DEFAULT_SKELETON_CONFIG_FILENAME,
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Invalid config file",
			arguments: arguments{
				configFileName: DEFAULT_INVALID_SKELETON_CONFIG_FILENAME,
			},
			wantError: true,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Unreadable config file",
			arguments: arguments{
				configFileName: DEFAULT_UNREADABLE_CONFIG_FILENAME,
			},
			wantError: true,
		},
	}

	fmt.Println(os.Getwd())

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				if _, errorInfo = ProcessBaseConfigFile(ts.arguments.configFileName); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Errorf(pi.UNEXPECTED_ERROR_FORMAT, tFunctionName, errorInfo.Error.Error())
				}
			},
		)
	}
}

func TestReadAndParseConfigFile(tPtr *testing.T) {

	type arguments struct {
		configFileName string
	}

	var (
		gotError           bool
		errorInfo          pi.ErrorInfo
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid config file",
			arguments: arguments{
				configFileName: fmt.Sprintf("%v%v", DEFAULT_SKELETON_CONFIG_FQD, DEFAULT_SKELETON_CONFIG_FILENAME),
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Unreadable config file",
			arguments: arguments{
				configFileName: fmt.Sprintf("%v%v", DEFAULT_SKELETON_CONFIG_FQD, DEFAULT_UNREADABLE_CONFIG_FILENAME),
			},
			wantError: true,
		},
	}

	fmt.Println(os.Getwd())

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				if _, errorInfo = ReadConfigFile(ts.arguments.configFileName); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Errorf(pi.UNEXPECTED_ERROR_FORMAT, tFunctionName, errorInfo.Error.Error())
				}
			},
		)
	}
}

func TestValidateConfiguration(tPtr *testing.T) {

	type arguments struct {
		config BaseConfiguration
	}

	var (
		gotError           bool
		errorInfo          pi.ErrorInfo
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "All missing",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  "",
					LogDirectory: "",
					MaxThreads:   0,
					PIDDirectory: "",
				},
			},
			wantError: true,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid settings Local Environment",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  ctv.ENVIRONMENT_LOCAL,
					LogDirectory: DEFAULT_LOG_DIRECTORY,
					MaxThreads:   DEFAULT_MAX_THREADS,
					PIDDirectory: DEFAULT_PID_DIRECTORY,
				},
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid settings Development Environment",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  ctv.ENVIRONMENT_DEVELOPMENT,
					LogDirectory: DEFAULT_LOG_DIRECTORY,
					MaxThreads:   DEFAULT_MAX_THREADS,
					PIDDirectory: DEFAULT_PID_DIRECTORY,
				},
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid settings Production Environment",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  ctv.ENVIRONMENT_PRODUCTION,
					LogDirectory: DEFAULT_LOG_DIRECTORY,
					MaxThreads:   DEFAULT_MAX_THREADS,
					PIDDirectory: DEFAULT_PID_DIRECTORY,
				},
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid settings Directories missing",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  ctv.ENVIRONMENT_PRODUCTION,
					LogDirectory: ctv.VAL_EMPTY,
					MaxThreads:   1,
					PIDDirectory: ctv.VAL_EMPTY,
				},
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid settings Negative threads",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  ctv.ENVIRONMENT_PRODUCTION,
					LogDirectory: ctv.VAL_EMPTY,
					MaxThreads:   -1,
					PIDDirectory: ctv.VAL_EMPTY,
				},
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid settings Zero threads",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  ctv.ENVIRONMENT_PRODUCTION,
					LogDirectory: ctv.VAL_EMPTY,
					MaxThreads:   0,
					PIDDirectory: ctv.VAL_EMPTY,
				},
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Valid settings Greater than threads cap",
			arguments: arguments{
				config: BaseConfiguration{
					Environment:  ctv.ENVIRONMENT_PRODUCTION,
					LogDirectory: ctv.VAL_EMPTY,
					MaxThreads:   THREAD_CAP + 1,
					PIDDirectory: ctv.VAL_EMPTY,
				},
			},
			wantError: false,
		},
	}

	fmt.Println(os.Getwd())

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				if errorInfo = ValidateConfiguration(ts.arguments.config); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Errorf(pi.UNEXPECTED_ERROR_FORMAT, tFunctionName, errorInfo.Error.Error())
				}
			},
		)
	}
}

// func TestValidateOptions(tPtr *testing.T) {
//
// 	type arguments struct {
// 		opts Configuration
// 	}
//
// 	//goland:noinspection ALL
// 	var (
// 		errorInfos         []pi.ErrorInfo
// 		gotError           bool
// 		tFunction, _, _, _ = runtime.Caller(0)
// 		tFunctionName      = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	tests := []struct {
// 		name      string
// 		arguments arguments
// 		wantError bool
// 	}{
// 		{
// 			name: "Positive Case: All options values are populated except for TLS and AuthenticatorService is Cognito.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated except for TLS and environment is local.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated except for TLS and environment is development.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_DEVELOPMENT,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated except for TLS and environment is production.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_PRODUCTION,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated and environment is production.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_PRODUCTION,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 					TLS: coreJWT.TLSInfo{
// 						TLSCert:     ctv.TEST_CERTIFICATE_FQN,
// 						TLSKey:      ctv.TEST_SAVUP_PRIVATE_KEY_FQN,
// 						TLSCABundle: ctv.TEST_CA_BUNDLE_FQN,
// 					},
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated except TLSCert and environment is production.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_PRODUCTION,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 					TLS: coreJWT.TLSInfo{
// 						TLSCert:     ctv.EMPTY,
// 						TLSKey:      ctv.TEST_SAVUP_PRIVATE_KEY_FQN,
// 						TLSCABundle: ctv.TEST_CA_BUNDLE_FQN,
// 					},
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated except TLSKey and environment is production.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_PRODUCTION,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 					TLS: coreJWT.TLSInfo{
// 						TLSCert:     ctv.TEST_CERTIFICATE_FQN,
// 						TLSKey:      ctv.EMPTY,
// 						TLSCABundle: ctv.TEST_CA_BUNDLE_FQN,
// 					},
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated except TLSCABundle and environment is production.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_PRODUCTION,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 					TLS: coreJWT.TLSInfo{
// 						TLSCert:     ctv.TEST_CERTIFICATE_FQN,
// 						TLSKey:      ctv.TEST_SAVUP_PRIVATE_KEY_FQN,
// 						TLSCABundle: ctv.EMPTY,
// 					},
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: All options values are populated except TLSCert and TLSKey and environment is production.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_PRODUCTION,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 					TLS: coreJWT.TLSInfo{
// 						TLSCert:     ctv.EMPTY,
// 						TLSKey:      ctv.EMPTY,
// 						TLSCABundle: ctv.EMPTY,
// 					},
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Negative Case: Authenticator Service is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.EMPTY,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: AWS credentials FQN is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.EMPTY,
// 					AWSInfoFQN:             ctv.TEST_NO_SUCH_FILE,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: All options values are populated and environment is missing.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.EMPTY,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: Missing Firebase Project Id",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.EMPTY,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: Firebase Credentials FQN is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_NO_SUCH_FILE,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: Firebase Credentials FQN has malformed JSON.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_MALFORMED_JSON_FILE,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: Log Directory is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_NO_SUCH_DIRECTORY,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: Log Directory is missing.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.EMPTY,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: PID Directory is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_NO_SUCH_DIRECTORY,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: PID Directory is missing",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.EMPTY,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Positive Case: Message Prefix is SAVUPLOCAL.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: Message Prefix is SAVUPDEV.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPDEV,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Positive Case: Message Prefix is SAVUP.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPPROD,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Negative Case: Message Prefix is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.TEST_STRING,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: NATS Creds FQN is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_NO_SUCH_FILE,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: NATS URL is missing.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.EMPTY,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: Private Key FQN is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_MALFORMED_JSON_FILE,
// 					SendGridKeyFQN:         ctv.TEST_GOOD_FQN,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: SendGrid Key FQN is invalid.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_NO_SUCH_FILE,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: SendGrid Key FQN has malformed JSON.",
// 			arguments: arguments{
// 				opts: Options{
// 					AuthenticatorService:   ctv.AUTH_COGNITO,
// 					AWSInfoFQN:             ctv.TEST_GOOD_FQN,
// 					Environment:            ctv.ENVIRONMENT_LOCAL,
// 					FirebaseProjectId:      ctv.TEST_STRING,
// 					FirebaseCredentialsFQN: ctv.TEST_GOOD_FQN,
// 					LogDirectory:           ctv.TEST_GOOD_FQN,
// 					PIDDirectory:           ctv.TEST_GOOD_FQN,
// 					MessagePrefix:          ctv.MESSAGE_PREFIX_SAVUPLOCAL,
// 					NATSCredsFQN:           ctv.TEST_GOOD_FQN,
// 					NATSURL:                ctv.TEST_STRING,
// 					PlaidKeyFQN:            ctv.TEST_GOOD_FQN,
// 					SendGridKeyFQN:         ctv.TEST_MALFORMED_JSON_FILE,
// 					StripeKeyFQN:           ctv.TEST_GOOD_FQN,
// 				},
// 			},
// 			wantError: true,
// 		},
// 	}
//
// 	for _, ts := range tests {
// 		tPtr.Run(ts.name, func(t *testing.T) {
// 			if errorInfos = validateOptions(ts.arguments.opts); len(errorInfos) > 0 {
// 				gotError = true
// 			} else {
// 				gotError = false
// 			}
// 			if gotError != ts.wantError {
// 				tPtr.Error(tFunctionName, ts.name, errorInfos)
// 			}
// 		})
// 	}
// }

// func TestCheckFileExistsAndReadable(tPtr *testing.T) {
//
// 	type arguments struct {
// 		FQN       string
// 		fileLabel string
// 	}
//
// 	var (
// 		errorInfo pi.ErrorInfo
// 		gotError  bool
// 	)
//
// 	tests := []struct {
// 		name      string
// 		arguments arguments
// 		wantError bool
// 	}{
// 		{
// 			name: "Positive Case: File exists and is readable.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_GOOD_FQN,
// 				fileLabel: "Test Good FQN",
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Negative Case: File doesn't exist.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_NO_SUCH_FILE,
// 				fileLabel: "Test Bad - No Such FQN",
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: File is not readable",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_UNREADABLE_FQN,
// 				fileLabel: "Test Bad - Unreadable FQN",
// 			},
// 			wantError: true,
// 		},
// 	}
//
// 	for _, ts := range tests {
// 		tPtr.Run(ts.name, func(t *testing.T) {
// 			if errorInfo = coreValidators.CheckFileExistsAndReadable(ts.arguments.FQN, ts.arguments.fileLabel); errorInfo.Error != nil {
// 				gotError = true
// 			} else {
// 				gotError = false
// 			}
// 			if gotError != ts.wantError {
// 				tPtr.Error(ts.name)
// 				tPtr.Error(errorInfo)
// 			}
// 		})
// 	}
// }

// func TestCheckFileValidJSON(tPtr *testing.T) {
//
// 	type arguments struct {
// 		FQN       string
// 		fileLabel string
// 	}
//
// 	var (
// 		errorInfo pi.ErrorInfo
// 		gotError  bool
// 	)
//
// 	tests := []struct {
// 		name      string
// 		arguments arguments
// 		wantError bool
// 	}{
// 		{
// 			name: "Positive Case: File contains valid JSON.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_GOOD_FQN,
// 				fileLabel: "Test Good JSON",
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Negative Case: File is not readable.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_UNREADABLE_FQN,
// 				fileLabel: "Test Unreadable File",
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: File contains INVALID JSON.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_MALFORMED_JSON_FILE,
// 				fileLabel: "Test Bad JSON",
// 			},
// 			wantError: true,
// 		},
// 	}
//
// 	for _, ts := range tests {
// 		tPtr.Run(ts.name, func(t *testing.T) {
// 			if errorInfo = coreValidators.CheckFileValidJSON(ts.arguments.FQN, ts.arguments.fileLabel); errorInfo.Error != nil {
// 				gotError = true
// 			} else {
// 				gotError = false
// 			}
// 			fmt.Println(gotError)
// 			if gotError != ts.wantError {
// 				tPtr.Error(ts.name)
// 				tPtr.Error(errorInfo)
// 			}
// 		})
// 	}
// }

// func TestReadAndParseConfigFile(tPtr *testing.T) {
//
// 	type arguments struct {
// 		FQN       string
// 		fileLabel string
// 	}
//
// 	var (
// 		errorInfo pi.ErrorInfo
// 		gotError  bool
// 	)
//
// 	tests := []struct {
// 		name      string
// 		arguments arguments
// 		wantError bool
// 	}{
// 		{
// 			name: "Positive Case: File contains valid JSON.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_GOOD_FQN,
// 				fileLabel: "Test Good JSON",
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Negative Case: File is not readable.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_UNREADABLE_FQN,
// 				fileLabel: "Test Unreadable File",
// 			},
// 			wantError: true,
// 		},
// 		{
// 			name: "Negative Case: File contains INVALID JSON.",
// 			arguments: arguments{
// 				FQN:       ctv.TEST_MALFORMED_JSON_FILE,
// 				fileLabel: "Test Bad JSON",
// 			},
// 			wantError: true,
// 		},
// 	}
//
// 	for _, ts := range tests {
// 		tPtr.Run(ts.name, func(t *testing.T) {
// 			if _, errorInfo = readAndParseConfigFile(ts.arguments.FQN); errorInfo.Error != nil {
// 				gotError = true
// 			} else {
// 				gotError = false
// 			}
// 			fmt.Println(gotError)
// 			if gotError != ts.wantError {
// 				tPtr.Error(ts.name)
// 				tPtr.Error(errorInfo)
// 			}
// 		})
// 	}
// }
