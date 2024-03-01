// Package sharedServices
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
package sharedServices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	// "net/httpServices"
	// "os"
	// "time"

	cj "GriesPikeThomp/shared-services/v2024/coreJWT"
	hv "GriesPikeThomp/shared-services/v2024/helpersValidators"
	ctv "github.com/sty-holdings/GriesPikeThomp/shared-services"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

type HTTPConfiguration struct {
	CredentialsFilename string      `json:"credentials_filename"`
	GinMode             string      `json:"gin_mode"`
	HTTPDomain          string      `json:"http_domain"`
	MessageEnvironment  string      `json:"message_environment"`
	Port                int         `json:"port"`
	RequestedThreads    uint        `json:"requested_threads"`
	RouteRegistry       []RouteInfo `json:"route_registry"`
	TLSInfo             cj.TLSInfo  `json:"tls_info"`
}

type RouteInfo struct {
	Namespace   string `json:"namespace"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type HTTPService struct {
	Config         HTTPConfiguration
	CredentialsFQN string
	HTTPServerPtr  *http.Server
	Secure         bool
}

// NewHTTP - creates a new httpServices service using the provided extension values.
//
//	Customer Messages: None
//	Errors: error returned by validateConfiguration
//	Verifications: validateConfiguration
func NewHTTP(configFilename string) (
	service HTTPService,
	errorInfo pi.ErrorInfo,
) {

	var (
		tAdditionalInfo = fmt.Sprintf("%v%v", ctv.TXT_FILENAME, configFilename)
		tConfig         HTTPConfiguration
		tConfigData     []byte
	)

	if tConfigData, errorInfo = config.ReadConfigFile(chv.PrependWorkingDirectory(configFilename)); errorInfo.Error != nil {
		return
	}

	if errorInfo.Error = json.Unmarshal(tConfigData, &tConfig); errorInfo.Error != nil {
		errorInfo = pi.NewErrorInfo(errorInfo.Error, tAdditionalInfo)
		return
	}

	if errorInfo = validateConfiguration(tConfig); errorInfo.Error != nil {
		return
	}

	service.Config = tConfig
	service.CredentialsFQN = hv.PrependWorkingDirectory(tConfig.CredentialsFilename)

	if tConfig.TLSInfo.TLSCert == ctv.VAL_EMPTY ||
		tConfig.TLSInfo.TLSPrivateKey == ctv.VAL_EMPTY ||
		tConfig.TLSInfo.TLSCABundle == ctv.VAL_EMPTY {
		service.Secure = false
	} else {
		service.Secure = true
	}

	return
}

//  Private Functions

// validateConfiguration - checks the NATS service configServers is valid.
//
//	Customer Messages: None
//	Errors: ErrEnvironmentInvalid, ErrMessageNamespaceInvalid, ErrDomainInvalid, error returned from DoesFileExistsAndReadable, ErrSubjectsMissing
//	Verifications: None
func validateConfiguration(config HTTPConfiguration) (errorInfo pi.ErrorInfo) {

	if errorInfo = hv.DoesFileExistsAndReadable(config.CredentialsFilename, ctv.TXT_FILENAME); errorInfo.Error != nil {
		pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, config.CredentialsFilename))
		return
	}
	if hv.IsBase64Encode(config.CredentialsFilename) == false {
		pi.NewErrorInfo(pi.ErrBase64Invalid, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, config.CredentialsFilename))
		return
	}
	if hv.IsGinModeValid(config.GinMode) == false {
		pi.NewErrorInfo(pi.ErrBase64Invalid, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, config.CredentialsFilename))
		return
	}
	if hv.IsEnvironmentValid(config.MessageEnvironment) == false {
		errorInfo = pi.NewErrorInfo(pi.ErrEnvironmentInvalid, fmt.Sprintf("%v%v", ctv.TXT_EVIRONMENT, config.MessageEnvironment))
		return
	}
	if hv.IsGinModeValid(config.GinMode) {
		config.GinMode = strings.ToLower(config.GinMode)
	} else {
		errorInfo = pi.NewErrorInfo(pi.ErrGinModeInvalid, fmt.Sprintf("%v%v", ctv.TXT_GIN_MODE, config.GinMode))
		return
	}
	if config.TLSInfo.TLSCert != ctv.VAL_EMPTY && config.TLSInfo.TLSPrivateKey != ctv.VAL_EMPTY && config.TLSInfo.TLSCABundle != ctv.VAL_EMPTY {
		if errorInfo = hv.DoesFileExistsAndReadable(config.TLSInfo.TLSCert, ctv.TXT_FILENAME); errorInfo.Error != nil {
			pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, config.TLSInfo.TLSCert))
			return
		}
		if errorInfo = hv.DoesFileExistsAndReadable(config.TLSInfo.TLSPrivateKey, ctv.TXT_FILENAME); errorInfo.Error != nil {
			pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, config.TLSInfo.TLSPrivateKey))
			return
		}
		if errorInfo = hv.DoesFileExistsAndReadable(config.TLSInfo.TLSCABundle, ctv.TXT_FILENAME); errorInfo.Error != nil {
			pi.NewErrorInfo(errorInfo.Error, fmt.Sprintf("%v%v", ctv.TXT_DIRECTORY, config.TLSInfo.TLSCABundle))
			return
		}
	}
	if len(config.RouteRegistry) == ctv.VAL_ZERO {
		pi.NewErrorInfo(pi.ErrSubjectsMissing, ctv.VAL_EMPTY)
	}

	return
}
