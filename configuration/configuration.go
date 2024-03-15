// Package sty_shared
/*
This is the STY-Holdings shared services

NOTES:

	Validation of the config must be done in the caller.

COPYRIGHT & WARRANTY:

	Copyright (c) 2022 STY-Holdings, inc
	All rights reserved.

	This software is the confidential and proprietary information of STY-Holdings, Inc.
	Use is subject to license terms.

	Unauthorized copying of this file, via any medium is strictly prohibited.

	Proprietary and confidential

	Written by Scott Yacko / syacko
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

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// GenerateConfigFileSkeleton will output to the console a skeleton file with notes.
//
//	Customer Messages: None
//	Errors: ErrConfigFileMissing
//	Verifications: None
func GenerateConfigFileSkeleton(serverName, SkeletonConfigFQD string) (errorInfo pi.ErrorInfo) {

	var (
		tSkeletonConfigData         []byte
		tSkeletonConfigNoteData     []byte
		tSkeletonConfigFilename     string
		tSkeletonConfigNoteFilename string
	)

	if serverName == ctv.VAL_EMPTY {
		pi.PrintError(pi.ErrServerNameMissing, fmt.Sprintf("%v %v", ctv.TXT_SERVER_NAME, serverName))
		return
	}
	if SkeletonConfigFQD == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrFileMissing, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, SkeletonConfigFQD))
		return
	}
	tSkeletonConfigFilename = fmt.Sprintf("%v%v", SkeletonConfigFQD, DEFAULT_SKELETON_CONFIG_FILENAME)
	tSkeletonConfigNoteFilename = fmt.Sprintf("%v%v", SkeletonConfigFQD, DEFAULT_SKELETON_CONFIG_NOTE_FILENAME)

	if tSkeletonConfigData, errorInfo.Error = os.ReadFile(tSkeletonConfigFilename); errorInfo.Error != nil {
		pi.PrintError(pi.ErrFileUnreadable, fmt.Sprintf("%v %v", ctv.TXT_FILENAME, tSkeletonConfigFilename))
		return
	}

	if tSkeletonConfigNoteData, errorInfo.Error = os.ReadFile(tSkeletonConfigNoteFilename); errorInfo.Error != nil {
		pi.PrintError(pi.ErrFileUnreadable, fmt.Sprintf("%v %v", ctv.TXT_FILENAME, tSkeletonConfigNoteFilename))
		return
	}

	fmt.Println("\nWhen '-g' is used all other program arguments are ignored.")
	fmt.Printf("\n%v Config file Skeleton: \n%v\n", serverName, string(tSkeletonConfigData))
	fmt.Println()
	fmt.Printf("%v\n", string(tSkeletonConfigNoteData))
	fmt.Println()

	return
}

// GetConfigFile - handles any configuration file.
//
//	Customer Messages: None
//	Errors: ReadConfigFile returned error, ErrJSONInvalid
//	Verifications: None
func GetConfigFile(
	configFileFQN string,
) (
	configData map[string]interface{},
	errorInfo pi.ErrorInfo,
) {

	var (
		tAdditionalInfo = fmt.Sprintf("%v%v", ctv.TXT_FILENAME, configFileFQN)
		tConfigData     []byte
		tConfigPtr      *map[string]interface{}
	)

	if tConfigData, errorInfo = ReadConfigFile(configFileFQN); errorInfo.Error != nil {
		return
	}

	if errorInfo.Error = json.Unmarshal(tConfigData, &tConfigPtr); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, tAdditionalInfo)
		return
	}

	configData = *tConfigPtr

	return
}

// ProcessBaseConfigFile - handles the base configuration file.
//
//	Customer Messages: None
//	Errors: errors returned from ReadConfigFile, ErrJSONInvalid
//	Verifications: None
func ProcessBaseConfigFile(configFileFQN string) (
	configData map[string]interface{},
	errorInfo pi.ErrorInfo,
) {

	if configData, errorInfo = GetConfigFile(configFileFQN); errorInfo.Error != nil {
		return
	}

	return
}

// ReadConfigFile opens the provide file, unmarshal the file and returns the Configuration object.
//
//	Customer Messages: None
//	Errors: ErrConfigFileMissing, ErrJSONInvalid
//	Verifications: None
func ReadConfigFile(configFileFQN string) (
	configData []byte,
	errorInfo pi.ErrorInfo,
) {

	var (
		tAdditionalInfo = fmt.Sprintf("%v %v", ctv.TXT_FILENAME, configFileFQN)
	)

	if configData, errorInfo.Error = os.ReadFile(configFileFQN); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(pi.ErrConfigFileMissing, tAdditionalInfo)
	}

	return
}

// Private function go below here
