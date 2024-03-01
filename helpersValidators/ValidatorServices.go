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
	"os"
	"reflect"
	"regexp"
	"strings"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// AreMapKeysPopulated - will test to make sure all map keys are set to anything other than nil or empty.
// func AreMapKeysPopulated(myMap map[any]interface{}) bool {
//
// 	if IsMapPopulated(myMap) {
// 		for key, _ := range myMap {
// 			if key == nil || key == ctv.TXT_EMPTY {
// 				return false
// 			}
// 		}
// 	} else {
// 		return false
// 	}
//
// 	return true
// }

// AreMapValuesPopulated - will test to make sure all map values are set to anything other than nil or empty.
// func AreMapValuesPopulated(myMap map[any]interface{}) bool {
//
// 	if IsMapPopulated(myMap) {
// 		for _, value := range myMap {
// 			if value == nil || value == ctv.VAL_EMPTY {
// 				return false
// 			}
// 		}
// 	} else {
// 		return false
// 	}
//
// 	return true
// }

// AreMapKeysValuesPopulated - check keys and value for missing values. Findings are ctv.GOOD, ctv.MISSING_VALUE,
// ctv.MISSING_KEY, or ctv.VAL_EMPTY_WORD.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: IsMapPopulated, AreMapKeysPopulated, AreMapValuesPopulated
// func AreMapKeysValuesPopulated(myMap map[any]interface{}) (finding string) {
//
// 	if IsMapPopulated(myMap) {
// 		if AreMapKeysPopulated(myMap) {
// 			if AreMapValuesPopulated(myMap) {
// 				finding = ctv.TXT_GOOD
// 			} else {
// 				finding = ctv.TXT_MISSING_VALUE
// 			}
// 		} else {
// 			finding = ctv.TXT_MISSING_KEY
// 		}
// 	} else {
// 		finding = ctv.TXT_EMPTY
// 	}
//
// 	return
// }

// DoesFileExistsAndReadable - works on any file. If the filename is not fully qualified
// the working directory will be prepended to the filename.
//
//	Customer Messages: None
//	Errors: ErrFileMissing, ErrFileUnreadable
//	Verifications: None
func DoesFileExistsAndReadable(filename, fileLabel string) (errorInfo pi.ErrorInfo) {

	var (
		fqn = PrependWorkingDirectory(filename)
	)

	if fileLabel == ctv.VAL_EMPTY {
		fileLabel = ctv.TXT_NO_LABEL_PROVIDED
	}
	errorInfo.AdditionalInfo = fmt.Sprintf("File: %v  Config File Label: %v", filename, fileLabel)

	if filename == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrFileMissing, errorInfo.AdditionalInfo)
		return
	}
	if DoesFileExist(fqn) == false {
		errorInfo = pi.NewErrorInfo(pi.ErrFileMissing, errorInfo.AdditionalInfo)
		return
	}
	if IsFileReadable(fqn) == false { // File is not readable
		errorInfo = pi.NewErrorInfo(pi.ErrFileUnreadable, errorInfo.AdditionalInfo)
	}

	return
}

// CheckFileValidJSON - reads the file and checks the contents
// func CheckFileValidJSON(FQN, fileLabel string) (errorInfo pi.ErrorInfo) {
//
// 	var (
// 		jsonData           []byte
// 		tFunction, _, _, _ = runtime.Caller(0)
// 		tFunctionName      = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
// 	errorInfo = pi.GetFunctionInfo()
// 	errorInfo.AdditionalInfo = fmt.Sprintf("File: %v  Config File Label: %v", FQN, fileLabel)
//
// 	if jsonData, errorInfo.Error = os.ReadFile(FQN); errorInfo.Error != nil {
// 		errorInfo.Error = pi.ErrFileUnreadable
// 		errorInfo.AdditionalInfo = fmt.Sprintf("FQN: %v File Label: %v", FQN, fileLabel)
// 		pi.PrintError(errorInfo)
// 	} else {
// 		if _isJSON := IsJSONValid(jsonData); _isJSON == false {
// 			errorInfo.Error = pi.ErrFileUnreadable
// 			errorInfo.AdditionalInfo = fmt.Sprintf("FQN: %v File Label: %v", FQN, fileLabel)
// 			pi.PrintError(errorInfo)
// 		}
// 	}
//
// 	return
// }

// DoesDirectoryExist - checks is the directory exists
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func DoesDirectoryExist(directoryName string) bool {

	return DoesFileExist(directoryName)
}

// DoesFileExist - does the value exist on the file system
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func DoesFileExist(fileName string) bool {

	if _, err := os.Stat(fileName); err == nil {
		return true
	}

	return false
}

// IsBase64Encode - will check if string is a valid base64 string.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsBase64Encode(base64Value string) bool {

	var (
		errorInfo pi.ErrorInfo
	)

	if _, errorInfo = Base64Decode(base64Value); errorInfo.Error == nil {
		return true
	}

	return false
}

// IsDirectoryFullyQualified - checks to see if the directory starts and ends with a slash.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsDirectoryFullyQualified(directory string) bool {

	if strings.HasPrefix(directory, ctv.FORWARD_SLASH) {
		if strings.HasSuffix(directory, ctv.FORWARD_SLASH) {
			return true
		}
	}

	return false

}

// IsDomainValid - checks if domain naming is followed
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsDomainValid(domain string) bool {

	if strings.ToLower(domain) == ctv.LOCAL_HOST {
		return true
	} else {
		regex := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
		if regex.MatchString(domain) {
			return true
		}
	}

	return false
}

// IsEnvironmentValid - checks that the value is valid. This function input is case-sensitive. Valid
// values are 'local', 'development', and 'production'.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsEnvironmentValid(environment string) bool {

	switch environment {
	case ctv.ENVIRONMENT_LOCAL:
	case ctv.ENVIRONMENT_DEVELOPMENT:
	case ctv.ENVIRONMENT_PRODUCTION:
	default:
		return false
	}

	return true
}

// IsEmpty - checks that the value is empty.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return isEmptyString(v)
	case []interface{}: // Check for slices and arrays
		return isEmptyCollection(v)
	case map[interface{}]interface{}: // Check for maps
		return isEmptyCollection(v)
	case chan interface{}: // Check for channels
		return isEmptyCollection(v)
	case *interface{}: // Check for pointers
		return isEmptyPointer(v)
	default:
		// For other types, consider if they have an "empty" equivalent
		return false
	}
}

// IsFileReadable - tries to open the file using 0644 permissions
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsFileReadable(fileName string) bool {

	if _, err := os.OpenFile(fileName, os.O_RDONLY, 0644); err == nil {
		return true
	}

	return false
}

// IsGinModeValid validates that the Gin httpServices framework mode is correctly set.
func IsGinModeValid(mode string) bool {

	switch strings.ToLower(mode) {
	case ctv.MODE_DEBUG:
	case ctv.MODE_RELEASE:
	default:
		return false
	}

	return true
}

// IsPopulated - checks that the value is populated.
func IsPopulated(value interface{}) bool {

	if IsEmpty(value) {
		return false
	}

	return true
}

// IsIPAddressValid - checks if the data provide is a valid IP address
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsIPAddressValid(ipAddress any) bool {

	// Checking if it is a valid IP addresses
	if IsIPv4Valid(ipAddress.(string)) || IsIPv6Valid(ipAddress.(string)) {
		return true
	}

	return false
}

// IsIPv4Valid - checks if the data provide is a valid IPv4 address
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsIPv4Valid(ipAddress any) bool {

	var (
		tIPv4Regex = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	)

	// Checking if it is a valid IPv4 addresses
	if tIPv4Regex.MatchString(ipAddress.(string)) {
		return true
	}

	return false
}

// IsIPv6Valid - checks if the data provide is a valid IPv6 address
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsIPv6Valid(ipAddress any) bool {

	var (
		tIPv6Regex = regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)
	)

	// Checking if it is a valid IPv4 addresses
	if tIPv6Regex.MatchString(ipAddress.(string)) {
		return true
	}

	return false
}

// IsJSONValid - checks if the data provide is valid JSON.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func IsJSONValid(jsonIn []byte) bool {

	var (
		jsonString map[string]interface{}
	)

	return json.Unmarshal(jsonIn, &jsonString) == nil
}

// IsMapPopulated - will determine if the map is populated.
// func IsMapPopulated(myMap map[any]interface{}) bool {
//
// 	if len(myMap) > 0 {
// 		return true
// 	}
//
// 	return false
// }

// IsMessagePrefixValid - is case-insensitive
// func IsMessagePrefixValid(messagePrefix string) bool {
//
// 	switch strings.ToUpper(messagePrefix) {
// 	case ctv.MESSAGE_PREFIX_SAVUPPROD:
// 	case ctv.MESSAGE_PREFIX_SAVUPDEV:
// 	case ctv.MESSAGE_PREFIX_SAVUPLOCAL:
// 	default:
// 		return false
// 	}
//
// 	return true
// }

// IsPeriodValid
// func IsPeriodValid(period string) bool {
//
// 	switch strings.ToUpper(period) {
// 	case ctv.YEAR:
// 	case ctv.MONTH:
// 	case ctv.DAY:
// 	default:
// 		return false
// 	}
//
// 	return true
// }

// This will set the connection values so GetConnection can be executed.
// func IsPostgresSSLModeValid(sslMode string) bool {
//
// 	switch sslMode {
// 	case ctv.POSTGRES_SSL_MODE_ALLOW:
// 	case ctv.POSTGRES_SSL_MODE_DISABLE:
// 	case ctv.POSTGRES_SSL_MODE_PREFER:
// 	case ctv.POSTGRES_SSL_MODE_REQUIRED:
// 	default:
// 		return false
// 	}
//
// 	return true
// }

// IsUserRegisterTypedValid
// func IsUserRegisterTypedValid(period string) bool {
//
// 	switch strings.ToUpper(period) {
// 	case ctv.COLLECTION_USER_TO_DO_LIST:
// 	case ctv.COLLECTION_USER_GOALS:
// 	default:
// 		return false
// 	}
//
// 	return true
// }

// IsURLValid
// func IsURLValid(URL string) bool {
//
// 	if _, err := url.ParseRequestURI(URL); err == nil {
// 		return true
// 	}
//
// 	return false
// }

// IsUUIDValid
// func IsUUIDValid(uuid string) bool {
//
// 	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
// 	return r.MatchString(uuid)
// }

// ValidateAuthenticatorService - Firebase is not support at this time
// func ValidateAuthenticatorService(authenticatorService string) (errorInfo pi.ErrorInfo) {
//
// 	var (
// 		tFunction, _, _, _ = runtime.Caller(0)
// 		tFunctionName      = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	switch strings.ToUpper(authenticatorService) {
// 	case ctv.AUTH_COGNITO:
// 	case ctv.AUTH_FIREBASE:
// 		fallthrough // ToDo This is because AUTH_FIREBASE is not supported right now
// 	default:
// 		errorInfo.Error = errors.New(fmt.Sprintf("The supplied authenticator service is not supported! Authenticator Service: %v (Authenticator Service is case insensitive)", authenticatorService))
// 		if authenticatorService == ctv.VAL_EMPTY {
// 			errorInfo.AdditionalInfo = "Authenticator Service parameter is empty"
// 		} else {
// 			errorInfo.AdditionalInfo = "Authenticator Service: " + authenticatorService
// 		}
// 	}
//
// 	return
// }

// ValidateDirectory - validates that the directory value is not empty and the value exists on the file system
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func ValidateDirectory(directory string) (errorInfo pi.ErrorInfo) {

	if directory == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, ctv.TXT_DIRECTORY_PARAM_EMPTY)
		return
	}
	if DoesDirectoryExist(directory) == false {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, directory))
	}

	return
}

// ValidateTransferMethod
// func ValidateTransferMethod(transferMethod string) (errorInfo pi.ErrorInfo) {
//
// 	var (
// 		tFunction, _, _, _ = runtime.Caller(0)
// 		tFunctionName      = runtime.FuncForPC(tFunction).Name()
// 	)
//
// 	pi.PrintDebugTrail(tFunctionName)
//
// 	switch strings.ToUpper(transferMethod) {
// 	case ctv.TRANFER_STRIPE:
// 	case ctv.TRANFER_WIRE:
// 	case ctv.TRANFER_CHECK:
// 	case ctv.TRANFER_ZELLE:
// 	default:
// 		errorInfo.Error = pi.ErrTransferMethodInvalid
// 		if transferMethod == ctv.VAL_EMPTY {
// 			errorInfo.AdditionalInfo = "Transfer Method parameter is empty"
// 		} else {
// 			errorInfo.AdditionalInfo = "Transfer Method: " + transferMethod
// 		}
// 	}
//
// 	return
// }

// Private methods below here

func isEmptyString(value string) bool {
	return value == ""
}

func isEmptyCollection(value interface{}) bool {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array || v.Kind() == reflect.Map || v.Kind() == reflect.Chan {
		return v.IsNil() || v.Len() == 0
	}
	return false
}

func isEmptyPointer(value interface{}) bool {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return false
}
