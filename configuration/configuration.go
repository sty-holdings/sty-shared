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
	"strings"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	hv "github.com/sty-holdings/sty-shared/v2024/helpersValidators"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// Configuration is a generic config file structure for application servers.
type BaseConfiguration struct {
	ConfigFQN         string
	DebugModeOn       bool                   `json:"debug_mode_on"`
	Environment       string                 `json:"environment"`
	LogDirectory      string                 `json:"log_directory"`
	MaxThreads        int                    `json:"max_threads"`
	PIDDirectory      string                 `json:"pid_directory"`
	Extensions        []BaseConfigExtensions `json:"load_extensions"`
	SkeletonConfigFQD string                 `json:"skeleton_config_fqd"`
}

type BaseConfigExtensions struct {
	Name           string `json:"name"`
	ConfigFilename string `json:"config_filename"`
}

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
		pi.PrintError(pi.ErrMissingServerName, fmt.Sprintf("%v %v", ctv.TXT_SERVER_NAME, serverName))
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
	configPtr interface{},
) (
	errorInfo pi.ErrorInfo,
) {

	var (
		tAdditionalInfo = fmt.Sprintf("%v %v", ctv.TXT_FILENAME, configFileFQN)
		tConfigData     []byte
	)

	if tConfigData, errorInfo = ReadConfigFile(configFileFQN); errorInfo.Error != nil {
		return
	}

	if errorInfo.Error = json.Unmarshal(tConfigData, &configPtr); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, tAdditionalInfo)
		return
	}

	return
}

// ProcessBaseConfigFile - handles the base configuration file.
//
//	Customer Messages: None
//	Errors: errors returned from ReadConfigFile, ErrJSONInvalid
//	Verifications: None
func ProcessBaseConfigFile(configFileFQN string) (
	config BaseConfiguration,
	errorInfo pi.ErrorInfo,
) {

	if errorInfo = GetConfigFile(configFileFQN, &config); errorInfo.Error != nil {
		return
	}

	config.ConfigFQN = configFileFQN
	config.Environment = strings.ToLower(config.Environment)

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

// ValidateConfiguration - checks the values in the configuration file are valid. ValidateConfiguration doesn't
// test if the configuration file exists, readable, or parsable. Defaults for LogDirectory, MaxThreads, and PIDDirectory
// are '/var/log/natsSerices-connect', 1, and '/var/run/natsSerices-connect', respectively.
//
//	Customer Messages: None
//	Errors: ErrEnvironmentInvalid, ErrDirectoryMissing, ErrMaxThreadsInvalid
//	Verifications: None
func ValidateConfiguration(config BaseConfiguration) (errorInfo pi.ErrorInfo) {

	if hv.IsEnvironmentValid(config.Environment) == false {
		errorInfo = pi.NewErrorInfo(pi.ErrEnvironmentInvalid, fmt.Sprintf("%v%v", ctv.TXT_EVIRONMENT, config.Environment))
		return
	}
	if hv.DoesDirectoryExist(config.LogDirectory) == false {
		pi.PrintError(pi.ErrDirectoryMissing, fmt.Sprintf("%v%v - Default Set: %v", ctv.TXT_DIRECTORY, config.LogDirectory, DEFAULT_LOG_DIRECTORY))
		config.LogDirectory = DEFAULT_LOG_DIRECTORY
	}
	if config.MaxThreads < 1 || config.MaxThreads > THREAD_CAP {
		pi.PrintError(pi.ErrMaxThreadsInvalid, fmt.Sprintf("%v%v - Default Set: %v", ctv.TXT_MAX_THREADS, config.LogDirectory, DEFAULT_MAX_THREADS))
		config.MaxThreads = DEFAULT_MAX_THREADS
	}
	if hv.DoesDirectoryExist(config.PIDDirectory) == false {
		pi.PrintError(pi.ErrDirectoryMissing, fmt.Sprintf("%v%v - Default Set: %v", ctv.TXT_DIRECTORY, config.LogDirectory, DEFAULT_PID_DIRECTORY))
		config.PIDDirectory = DEFAULT_PID_DIRECTORY
	}
	if hv.DoesDirectoryExist(config.SkeletonConfigFQD) == false {
		pi.PrintError(pi.ErrDirectoryMissing, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, config.SkeletonConfigFQD))
		config.LogDirectory = DEFAULT_LOG_DIRECTORY
	}

	return
}

// Private function go below here
