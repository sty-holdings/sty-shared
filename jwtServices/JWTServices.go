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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	hv "github.com/sty-holdings/sty-shared/v2024/helpersValidators"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// BuildTLSTemporaryFiles - creates temporary files for TLS information.
// The function checks if the TLSCABundle, TLSCert, and TLSPrivateKey in tlsInfo are provided. If any of these values are empty,
// the function returns an error indicating the missing.
//
//	Customer Messages: None
//	Errors: ErrRequiredArgumentMissing, returned from WriteFile
//	Verifications: None
func BuildTLSTemporaryFiles(
	tempDirectory string,
	tlsInfo TLSInfo,
) (
	tlsFQN map[string]string,
	errorInfo pi.ErrorInfo,
) {

	if tlsInfo.TLSCABundle == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_MISSING_PARAMETER, ctv.FN_TLS_CA_BUNDLE))
		return
	} else {
		if errorInfo = hv.WriteFile(fmt.Sprintf("%v/%v", tempDirectory, TLS_CA_BUNDLE_FILENAME), []byte(tlsInfo.TLSCABundle), 0744); errorInfo.Error != nil {
			return
		}
	}
	if tlsInfo.TLSCert == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_MISSING_PARAMETER, ctv.FN_TLS_CERTIFICATE))
		return
	} else {
		if errorInfo = hv.WriteFile(fmt.Sprintf("%v/t%v", tempDirectory, TLS_CERT_FILENAME), []byte(tlsInfo.TLSCert), 0744); errorInfo.Error != nil {
			return
		}
	}
	if tlsInfo.TLSPrivateKey == ctv.VAL_EMPTY {
		errorInfo = pi.NewErrorInfo(pi.ErrRequiredArgumentMissing, fmt.Sprintf("%v%v", ctv.TXT_MISSING_PARAMETER, ctv.FN_TLS_PRIVATE_KEY))
		return
	} else {
		if errorInfo = hv.WriteFile(fmt.Sprintf("%v/%v", tempDirectory, TLS_PRIVATE_KEY), []byte(tlsInfo.TLSPrivateKey), 0744); errorInfo.Error != nil {
			return
		}
	}

	return
}

// GenerateJWT
// Create a new token object, specifying signing method and the claims
// you would like it to contain.
// func GenerateJWT(privateKey, requestorId, period string, duration int64) (jwtServices string, errorInfo pi.ErrorInfo) {
//
// 	var (
// 		tDuration      time.Duration
// 		tPrivateKey    *rsa.PrivateKey
// 		tRawPrivateKey []byte
// 	)
//
// 	if privateKey == ctv.VAL_EMPTY {
// 		errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! %v: '%v'", ctv.FN_PRIVATE_KEY, ctv.VAL_EMPTY))
// 		log.Println(errorInfo.Error)
// 	} else {
// 		if requestorId == ctv.VAL_EMPTY || period == ctv.VAL_EMPTY || duration < 1 {
// 			errorInfo.Error = errors.New(fmt.Sprintf("Require information is missing! %v: '%v' %v: '%v' %v: '%v'", ctv.FN_REQUESTOR_ID, requestorId, ctv.FN_PERIOD, period, ctv.FN_DURATION, duration))
// 			log.Println(errorInfo.Error)
// 		} else {
// 			if cv.IsPeriodValid(period) && duration > 0 {
// 				tRawPrivateKey = []byte(privateKey)
// 				if tPrivateKey, errorInfo = ParsePrivateKey(tRawPrivateKey); errorInfo.Error == nil {
// 					switch strings.ToUpper(period) {
// 					case "M":
// 						tDuration = time.Minute * time.Duration(duration)
// 					case "H":
// 						tDuration = time.Hour * time.Duration(duration)
// 					case "D":
// 						tDuration = time.Hour * time.Duration(duration*24)
// 					default:
// 						tDuration = time.Hour * time.Duration(duration)
// 					}
// 					jwtServices, errorInfo.Error = jwt2.NewWithClaims(jwt2.SigningMethodRS512, jwt2.MapClaims{
// 						"requestorId": requestorId,
// 						"Issuer":      ctv.CERT_ISSUER,
// 						"Subject":     requestorId,
// 						"ExpiresAt":   time.Now().Add(tDuration).String(),
// 						"NotBefore":   time.Now(),
// 					}).SignedString(tPrivateKey)
// 				}
// 			}
// 		}
// 	}
//
// 	return
// }

// GenerateRSAKey
func GenerateRSAKey(rsaBits int) (
	privateKey crypto.PrivateKey,
	publicKey crypto.PublicKey,
	errorInfo pi.ErrorInfo,
) {

	var (
		_PrivateKey *rsa.PrivateKey
	)

	if _PrivateKey, errorInfo.Error = rsa.GenerateKey(rand.Reader, rsaBits); errorInfo.Error != nil {
		log.Println(errorInfo.Error)
	}

	if errorInfo.Error == nil {
		// The public key is a part of the *rsa.PrivateKey struct
		publicKey = _PrivateKey.Public()
		privateKey = _PrivateKey
	}

	return
}

// ParsePrivateKey
func ParsePrivateKey(tRawPrivateKey []byte) (
	privateKey *rsa.PrivateKey,
	errorInfo pi.ErrorInfo,
) {

	if privateKey, errorInfo.Error = jwt.ParseRSAPrivateKeyFromPEM(tRawPrivateKey); errorInfo.Error != nil {
		errorInfo.Error = errors.New("Unable to parse the private key referred to in the configuration file.")
		log.Println(errorInfo.Error)
	}

	return
}

// RemoveTLSTemporaryFiles - removes the temporary CA Bundle, Certificate, and Private Key files.
//
//	Customer Messages: None
//	Errors: Return from RemoveFile
//	Verifications: None
func RemoveTLSTemporaryFiles(
	tempDirectory string,
) (errorInfo pi.ErrorInfo) {

	if errorInfo = hv.RemoveFile(fmt.Sprintf("%v/tls-ca-bundle.crt", tempDirectory)); errorInfo.Error == nil {
		if errorInfo = hv.RemoveFile(fmt.Sprintf("%v/tls-ca-cert.crt", tempDirectory)); errorInfo.Error == nil {
			if errorInfo = hv.RemoveFile(fmt.Sprintf("%v/tls-private.key", tempDirectory)); errorInfo.Error == nil {
			}
		}
	}

	return
}
