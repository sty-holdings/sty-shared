// Package sty-shared
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
	"crypto"
)

//goland:noinspection GoSnakeCaseUsage,GoCommentStart
const (
	TLS_CA_BUNDLE_FILENAME   = "tls-ca-bundle.crt"
	TLS_CERT_FILENAME        = "tls-cert.crt"
	TLS_PRIVATE_KEY_FILENAME = "tls-private.key"
)

type GenerateCertificate struct {
	CertFileName       string
	Certificate        []byte
	Host               string
	PublicKey          crypto.PublicKey
	PrivateKey         crypto.PrivateKey
	PrivateKeyFileName string
	RSABits            int
	SelfCA             bool
	ValidFor           string
}

// TLSInfo files
type TLSInfo struct {
	TLSCert          string `json:"tls_certificate"`
	TLSCertFQN       string `json:"tls_certificate_fqn"`
	TLSPrivateKey    string `json:"tls_private_key"`
	TLSPrivateKeyFQN string `json:"tls_private_key_fqn"`
	TLSCABundle      string `json:"tls_ca_bundle"`
	TLSCABundleFQN   string `json:"tls_ca_bundle_fqn"`
}
