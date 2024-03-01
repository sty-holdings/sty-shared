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
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	rcv "github.com/sty-holdings/resuable-const-vars/src"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// Base64Decode - will decode a base64 string to a string. If there is an error,
// the first 20 characters of the base64 string are logged.
// REMINDER: If the base64 string has sensitivity information, empty out the
// ErrorInfo.AdditionalInfo field before logging or outputting the error.
//
//	Customer Messages: None
//	Errors: error returned by StdEncoding.DecodeString
//	Verifications: None
func Base64Decode(base64Value string) (
	value []byte,
	errorInfo pi.ErrorInfo,
) {

	if value, errorInfo.Error = b64.StdEncoding.DecodeString(base64Value); errorInfo.Error != nil {
		errorInfo.AdditionalInfo = fmt.Sprintf("%v%v", ctv.TXT_BASE64, base64Value[:20])
	}

	return
}

// Base64Encode - will encode a string to a base64 string
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func Base64Encode(value string) string {

	return b64.StdEncoding.EncodeToString([]byte(value))
}

// BuildJSONRequest
// func BuildJSONRequest(request interface{}) (jsonRequest []byte) {
//
// 	var (
// 		err error
// 	)
//
// 	if jsonRequest, err = json.Marshal(request); err != nil {
// 		err = errors.New(fmt.Sprintf("Failed to generate JSON payload. Error: %v", err.Error()))
// 		log.Println(err.Error())
// 		// 	todo Error Handling & Notification
// 	}
//
// 	if coreValidators.IsJSONValid(jsonRequest) == false {
// 		jsonRequest = nil
// 		err = errors.New(fmt.Sprintf("Was not able to generate valid json for request %v", request))
// 		log.Println(err.Error())
// 		// 	todo Error Handling & Notification
// 	}
//
// 	return
// }

// BuildLegalName
// func BuildLegalName(firstName, lastName string) (legalName string) {
//
// 	if firstName != ctv.EMPTY && lastName != ctv.EMPTY {
// 		legalName = fmt.Sprintf("%v %v", firstName, lastName)
// 	}
//
// 	return
// }

// CapitalizeFirstLetter - will make the first letter of the string to upper case and the other letters to lower
// func CapitalizeFirstLetter(stringIn string) string {
//
// 	if stringIn == ctv.EMPTY {
// 		return ctv.EMPTY
// 	}
//
// 	x := []byte(stringIn)
// 	y := bytes.ToUpper([]byte{x[0]})
// 	z := bytes.ToLower(x[1:])
//
// 	return string(bytes.Join([][]byte{y, z}, nil))
// }

// ConvertMapAnyToMapString
// func ConvertMapAnyToMapString(mapIn map[any]interface{}) (mapOut map[string]interface{}) {
//
// 	mapOut = make(map[string]interface{})
//
// 	if coreValidators.IsMapPopulated(mapIn) {
// 		for key, value := range mapIn {
// 			mapOut[key.(string)] = value
// 		}
// 	}
//
// 	return
// }

// ConvertSliceToSliceOfPtrs - takes a slice and returns a slice of pointers to the items in the slice.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func ConvertStringSliceToSliceOfPtrs(inbound []string) (outbound []*string) {

	for _, tSlice := range inbound {
		x := tSlice // required to get different pointers for each value.
		outbound = append(outbound, &x)
	}

	return
}

// CreateAndRedirectLogOutput - will create the fully qualified config file log directory.
// The log output is based on the redirectTo value, [MODE_OUTPUT_LOG | MODE_OUTPUT_LOG_DISPLAY].
// The log file name uses this format: 2006-01-02 15:04:05.000 Z0700. All spaces, colons, and periods
// are replaced with underscores.
//
//	Customer Messages: None
//	Errors: ErrDirectoryNotFullyQualified, any error from os.OpenFile
//	Verifications: IsDirectoryFullyQualified
func CreateAndRedirectLogOutput(logDirectory, redirectTo string) (
	logFileHandlerPtr *os.File,
	logFQN string,
	errorInfo pi.ErrorInfo,
) {

	switch redirectTo {
	case ctv.MODE_OUTPUT_LOG:
		logFileHandlerPtr, logFQN, errorInfo = createLogFile(logDirectory)
		log.SetOutput(io.MultiWriter(logFileHandlerPtr))
	case ctv.MODE_OUTPUT_LOG_DISPLAY:
		logFileHandlerPtr, logFQN, errorInfo = createLogFile(logDirectory)
		log.SetOutput(io.MultiWriter(os.Stdout, logFileHandlerPtr))
	default:
		errorInfo = pi.NewErrorInfo(pi.ErrMissingServerName, fmt.Sprintf("%v%v", ctv.TXT_REDIRECT, redirectTo))
	}

	return
}

// DoesFieldExist - tests the struct for the field name.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func DoesFieldExist(
	structType interface{},
	fieldName string,
) bool {

	var (
		found bool
	)

	_, found = reflect.TypeOf(structType).FieldByName(fieldName)

	return found
}

// FloatToPennies - multiples the value by 100. Called pennies because we did for the US first.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func FloatToPennies(amount float64) (pennies int64) {

	return int64(amount * 100)
}

// formatURL - will return a formatted url with the protocol, domain, and port.
//
//	Validation: none
//	Format: "%v://%v:%v"
//	Example: http://verifyemail.savup.com:2134, https://verifyemail.savup.com:2134, http://localhost:2134, https://localhost:2134
// func formatURL(protocol, domain string, port uint) (url string) {
//
// 	if domain == ctv.ENVIRONMENT_LOCAL {
// 		url = fmt.Sprintf("%v://%v:%v", protocol, ctv.HTTP_DOMAIN_LOCALHOST, port)
// 	} else {
// 		url = fmt.Sprintf("%v://%v:%v", protocol, domain, port)
// 	}
//
// 	return
// }

// GenerateEndDate - will return a string by taking the startDate and adding months.
// If the startDate is empty the endDate will be empty.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
// func GenerateEndDate(startDate string, months int) (endDate string) {
//
// 	var (
// 		err    error
// 		tStart time.Time
// 	)
//
// 	if startDate == "" {
// 		endDate = ""
// 	} else {
// 		// Parse the start date string.
// 		tStart, err = time.Parse("2006-01-02", startDate)
// 		if err != nil {
// 			panic(err)
// 		}
// 		// Calculate the end date.
// 		end := tStart.AddDate(0, months, 0)
// 		endDate = end.Format("2006-01-02")
// 	}
//
// 	return
// }

// GenerateURL - will return the protocol, domain, and port. Using HTTP_PROTOCOL_SECURE or HTTP_PROTOCOL_NON_SECURE,
// ENDPOINT_VERIFY_EMAIL and HTTP_SECURE_PORT or HTTP_NON_SECURE_PORT based on the arguments.
//
//	Customer Messages: None
//	Errors: None
//	Verification: none
//	Example: http://localhost:1234, https://localhost:1234, http://api-dev.savup.com:1234, https://api-dev.savup.com:1234
//
// ToDo Change the Environment_local domain to local host once we have resolved the handshake issue happening when savup-httpServices is run locally.
// func GenerateURL(environment string, secure bool) (url string) {
//
// 	switch fmt.Sprintf("%v,%v", strings.ToUpper(environment), secure) {
// 	case fmt.Sprintf("%v,%v", ctv.ENVIRONMENT_LOCAL, true):
// 		url = formatURL(ctv.HTTP_PROTOCOL_SECURE, ctv.HTTP_DOMAIN_API_LOCAL, ctv.HTTP_PORT_SECURE)
// 	case fmt.Sprintf("%v,%v", ctv.ENVIRONMENT_LOCAL, false):
// 		url = formatURL(ctv.HTTP_PROTOCOL_NON_SECURE, ctv.HTTP_DOMAIN_API_LOCAL, ctv.HTTP_PORT_NON_SECURE)
// 	case fmt.Sprintf("%v,%v", ctv.ENVIRONMENT_DEVELOPMENT, true):
// 		url = formatURL(ctv.HTTP_PROTOCOL_SECURE, ctv.HTTP_DOMAIN_API_DEV, ctv.HTTP_PORT_SECURE)
// 	case fmt.Sprintf("%v,%v", ctv.ENVIRONMENT_DEVELOPMENT, false):
// 		url = formatURL(ctv.HTTP_PROTOCOL_NON_SECURE, ctv.HTTP_DOMAIN_API_DEV, ctv.HTTP_PORT_NON_SECURE)
// 	case fmt.Sprintf("%v,%v", ctv.ENVIRONMENT_PRODUCTION, true):
// 		url = formatURL(ctv.HTTP_PROTOCOL_SECURE, ctv.HTTP_DOMAIN_API_PROD, ctv.HTTP_PORT_SECURE)
// 	case fmt.Sprintf("%v,%v", ctv.ENVIRONMENT_PRODUCTION, false):
// 		url = formatURL(ctv.HTTP_PROTOCOL_NON_SECURE, ctv.HTTP_DOMAIN_API_PROD, ctv.HTTP_PORT_NON_SECURE)
// 	}
//
// 	return
// }

// GenerateUUIDType1
// func GenerateUUIDType1(removeDashes bool) (myUUID string) {
//
// 	_UUID, _ := uuid.NewUUID()
// 	myUUID = fmt.Sprintf("%v", _UUID)
//
// 	if removeDashes {
// 		myUUID = strings.Replace(myUUID, "-", "", -1)
// 	}
//
// 	return
// }

// GenerateUUIDType4
// func GenerateUUIDType4(removeDashes bool) (myUUID string) {
//
// 	_UUID, _ := uuid.NewRandom()
// 	myUUID = fmt.Sprintf("%v", _UUID)
//
// 	if removeDashes {
// 		myUUID = strings.Replace(myUUID, "-", "", -1)
// 	}
//
// 	return
// }

// GenerateVerifyEmailURLWithUUID - return the url and uuid for the Verify Email.
// func GenerateVerifyEmailURLWithUUID(environment string, secure bool) (url, uuid string) {
//
// 	uuid = GenerateUUIDType4(false)
// 	url = fmt.Sprintf("%v?uuid=%v", GenerateVerifyEmailURL(environment, secure), uuid)
//
// 	return
// }

// GenerateVerifyEmailURLWithUUIDUsername - return the url, uuid and the username for the Verify Email.
// func GenerateVerifyEmailURLWithUUIDUsername(username, environment string, secure bool) (url, uuid string) {
//
// 	uuid = GenerateUUIDType4(false)
// 	url = fmt.Sprintf("%v?%v=%v&%v=%v", GenerateVerifyEmailURL(environment, secure), ctv.FN_UUID, uuid, ctv.FN_USERNAME, username)
//
// 	return
// }

// GenerateVerifyEmailURL - return the url.
// func GenerateVerifyEmailURL(environment string, secure bool) (url string) {
//
// 	url = fmt.Sprintf("%v/%v", GenerateURL(environment, secure), ctv.ENDPOINT_VERIFY_EMAIL)
//
// 	return
// }

// GetDate - return the current date in YYYY-MM-DD format
//
//	Customer Message: None
//	Errors: None
//	Verification: None
// func GetDate() string {
// 	return time.Now().Format("2006-01-02")
// }

// GetTime - return the current time in HH-mm-ss.00000 format, where hour is in military time.
//
//	Customer Message: None
//	Errors: None
//	Verification: None
// func GetTime() string {
// 	return time.Now().Format("15-04-05.00000")
// }

// GetFieldsNames - will return a list fields in a struct
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func GetFieldsNames(unknownStruct interface{}) (
	fields map[string]interface{},
	errorInfo pi.ErrorInfo,
) {

	fields = make(map[string]interface{})

	tStruct := reflect.ValueOf(unknownStruct)
	tType := tStruct.Type()

	for i := 0; i < tType.NumField(); i++ {
		if tType.Field(i).IsExported() {
			fields[tType.Field(i).Name] = tStruct.FieldByName(tType.Field(i).Name).Interface()
		}
	}

	return
}

// PenniesToFloat
// func PenniesToFloat(pennies int64) float64 {
//
// 	return float64(pennies) / 100
// }

// PrintAndDie - is exported for access in other packages. Not going to test
// func PrintAndDie(msg string) {
//
// 	_, _ = fmt.Fprintln(os.Stderr, msg)
// 	os.Exit(1)
//
// }

// PrependWorkingDirectory - will add the working directory.
// if the filename first character is a /, the passed value will be returned
// unmodified.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func PrependWorkingDirectory(filename string) string {

	var (
		tWorkingDirectory, _ = os.Getwd()
	)

	if filepath.IsAbs(filename) {
		return filename
	}

	return fmt.Sprintf("%v/%v", tWorkingDirectory, filename)
}

// PrependWorkingDirectoryWithEndingSlash - will add the working directory, a slash, the directory
// provided, and an ending slash. If the directory first character is a slash, the passed value will
// be returned unmodified. The last character is not checked, so you could end up with two slashes.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func PrependWorkingDirectoryWithEndingSlash(directory string) string {

	var (
		tWorkingDirectory, _ = os.Getwd()
	)

	if filepath.IsAbs(directory) {
		return directory
	}

	return fmt.Sprintf("%v/%v/", tWorkingDirectory, directory)
}

// printDashLine - will output a given number of dashed lines based on the outputMode.
// The default is to output to the log
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
// func printDashLines(lines int, outputMode string) {
//
// 	for i := 0; i < lines; i++ {
// 		if strings.ToLower(outputMode) == ctv.MODE_OUTPUT_DISPLAY {
// 			fmt.Println("------------------------------------------")
// 		} else {
// 			log.Println("------------------------------------------")
// 		}
// 	}
// }

// PrintLinesAtStartOfRequest - will output dashed lines when a new request arrives.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
// func PrintLinesAtStartOfRequest(lines int, outputMode string) {
//
// 	printDashLines(lines, outputMode)
//
// }

// RedirectLogOutput - will redirect log output based on the redirectTo value, [MODE_OUTPUT_LOG | MODE_OUTPUT_LOG_DISPLAY].
//
//	Customer Messages: None
//	Errors: ErrDirectoryNotFullyQualified, any error from os.OpenFile
//	Verifications: IsDirectoryFullyQualified
func RedirectLogOutput(
	inLogFileHandlerPtr *os.File,
	redirectTo string,
) (errorInfo pi.ErrorInfo) {

	switch redirectTo {
	case ctv.MODE_OUTPUT_LOG:
		log.SetOutput(io.MultiWriter(inLogFileHandlerPtr))
	case ctv.MODE_OUTPUT_LOG_DISPLAY:
		log.SetOutput(io.MultiWriter(os.Stdout, inLogFileHandlerPtr))
	default:
		errorInfo = pi.NewErrorInfo(pi.ErrRedirectModeInvalid, fmt.Sprintf("%v%v", ctv.TXT_REDIRECT, redirectTo))
	}

	return
}

// RemovePidFile - removes the pid file for the running instance
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func RemovePidFile(pidFQN string) (errorInfo pi.ErrorInfo) {

	// This doesn't use the coreValidator.DoesFileExist by design.
	if _, err := os.Stat(pidFQN); err != nil {
		errorInfo = pi.NewErrorInfo(pi.ErrFileMissing, fmt.Sprintf("%v%v", ctv.TXT_FILENAME, pidFQN))
		return
	}

	if errorInfo.Error = os.Remove(pidFQN); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(pi.ErrFileRemovalFailed, fmt.Sprintf("%v%v", ctv.TXT_FILENAME, pidFQN))
		return
	}

	return
}

// WriteFile - will create and write to a fully qualified file.
//
//	Customer Messages: None
//	Errors: ErrFileCreationFailed
//	Verifications: None
func WriteFile(
	fqn string,
	fileData []byte,
	filePermissions os.FileMode,
) (errorInfo pi.ErrorInfo) {

	if errorInfo.Error = os.WriteFile(fqn, fileData, filePermissions); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v %v%v", pi.ErrFileCreationFailed.Error(), ctv.TXT_FILENAME, fqn))
	}

	return
}

// WritePidFile - will create and write the server pid file.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func WritePidFile(
	pidFQN string,
	pid int,
) (errorInfo pi.ErrorInfo) {

	if errorInfo = WriteFile(pidFQN, []byte(strconv.Itoa(pid)), 0766); errorInfo.Error == nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_FILENAME, pidFQN))
	}

	return
}

// Private Functions

// createLogFile - will create and open the  log file using the fully qualified directory.
//
//	Customer Messages: None
//	Errors: ErrDirectoryNotFullyQualified, any error from os.OpenFile
//	Verifications: IsDirectoryFullyQualified
func createLogFile(logFQD string) (
	logFileHandlerPtr *os.File,
	logFQN string,
	errorInfo pi.ErrorInfo,
) {

	var (
		tLogFileName string
	)

	if IsDirectoryFullyQualified(logFQD) == false {
		errorInfo = pi.NewErrorInfo(pi.ErrDirectoryNotFullyQualified, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, logFQD))
		return
	}

	tDateTime := time.Now().Format("2006-01-02 15:04:05.000 Z0700")
	tLogFileName = strings.Replace(
		strings.Replace(strings.Replace(tDateTime, ctv.SPACES_ONE, ctv.UNDERSCORE, -1), ctv.COLON, ctv.UNDERSCORE, -1),
		ctv.PERIOD,
		ctv.UNDERSCORE,
		-1,
	)
	logFQN = fmt.Sprintf("%v%v.log", logFQD, tLogFileName)

	// Set log file output
	if logFileHandlerPtr, errorInfo.Error = os.OpenFile(logFQN, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_FILENAME, logFQN))
		return
	}

	return
}

// getType
// func getType(myVar interface{}) (myType string) {
//
// 	if t := reflect.TypeOf(myVar); t.Kind() == reflect.Ptr {
// 		myType = "*" + t.Elem().Name()
// 	} else {
// 		myType = t.Name()
// 	}
//
// 	return
// }
