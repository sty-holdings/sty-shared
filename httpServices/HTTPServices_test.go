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

	ctv "github.com/sty-holdings/GriesPikeThomp/shared-services"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

func TestNewHTTP(tPtr *testing.T) {

	type arguments struct {
		hostname       string
		configFilename string
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
			name: ctv.TEST_POSITIVE_SUCCESS + "Secure connection.",
			arguments: arguments{
				hostname:       "localhost",
				configFilename: "/Users/syacko/workspace/sty-holdings/GriesPikeThomp/servers/natsSerices-connect/config/local/httpServices-inbound-config.json",
			},
			wantError: false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Bad URL.",
			arguments: arguments{
				hostname:       "localhost",
				configFilename: "/Users/syacko/workspace/sty-holdings/GriesPikeThomp/servers/natsSerices-connect/config/local/httpServices-inbound-config.json",
			},
			wantError: true,
		},
		{
			name: ctv.TEST_NEGATIVE_SUCCESS + "Missing credentials location.",
			arguments: arguments{
				hostname:       "localhost",
				configFilename: "/Users/syacko/workspace/sty-holdings/GriesPikeThomp/servers/natsSerices-connect/config/local/httpServices-inbound-config.json",
			},
			wantError: true,
		},
		{
			name: ctv.TEST_NEGATIVE_SUCCESS + "Missing certificate FQN.",
			arguments: arguments{
				hostname:       "localhost",
				configFilename: "/Users/syacko/workspace/sty-holdings/GriesPikeThomp/servers/natsSerices-connect/config/local/httpServices-inbound-config.json",
			},
			wantError: true,
		},
		{
			name: ctv.TEST_NEGATIVE_SUCCESS + "Missing private key FQN.",
			arguments: arguments{
				hostname:       "localhost",
				configFilename: "/Users/syacko/workspace/sty-holdings/GriesPikeThomp/servers/natsSerices-connect/config/local/httpServices-inbound-config.json",
			},
			wantError: true,
		},
		{
			name: ctv.TEST_NEGATIVE_SUCCESS + "Missing CA bundle FQN.",
			arguments: arguments{
				hostname:       "localhost",
				configFilename: "/Users/syacko/workspace/sty-holdings/GriesPikeThomp/servers/natsSerices-connect/config/local/httpServices-inbound-config.json",
			},
			wantError: true,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				if _, errorInfo = NewHTTP(ts.arguments.configFilename); errorInfo.Error != nil {
					gotError = true
				} else {
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Error(ts.name)
					tPtr.Error(errorInfo)
				}
			},
		)
	}
}
