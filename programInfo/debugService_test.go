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
	"runtime"
	"testing"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
)

func TestPrintDebugFunctionInfo(tPtr *testing.T) {

	type arguments struct {
		debugModeOn bool
		outputMode  string
	}

	var (
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: "Positive Case: Debug False - No output.",
			arguments: arguments{
				debugModeOn: false,
				outputMode:  "",
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug true - No output.",
			arguments: arguments{
				debugModeOn: true,
				outputMode:  "",
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug False - Output Display.",
			arguments: arguments{
				debugModeOn: false,
				outputMode:  ctv.MODE_OUTPUT_DISPLAY,
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug True - Output Display.",
			arguments: arguments{
				debugModeOn: true,
				outputMode:  ctv.MODE_OUTPUT_DISPLAY,
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug False - Output Log.",
			arguments: arguments{
				debugModeOn: false,
				outputMode:  ctv.MODE_OUTPUT_LOG,
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug True - Output Log.",
			arguments: arguments{
				debugModeOn: true,
				outputMode:  ctv.MODE_OUTPUT_LOG,
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			tFunctionName, func(t *testing.T) {
				PrintDebugFunctionInfo(ts.arguments.debugModeOn, ts.arguments.outputMode)
			},
		)
	}
}

func TestPrintDebugLine(tPtr *testing.T) {

	type arguments struct {
		message     string
		debugModeOn bool
		outputMode  string
	}

	var (
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
	)

	tests := []struct {
		name      string
		arguments arguments
		wantError bool
	}{
		{
			name: "Positive Case: Debug False - No output.",
			arguments: arguments{
				debugModeOn: false,
				outputMode:  "",
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug true - No output.",
			arguments: arguments{
				debugModeOn: true,
				outputMode:  "",
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug False - Output Display.",
			arguments: arguments{
				debugModeOn: false,
				outputMode:  ctv.MODE_OUTPUT_DISPLAY,
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug True - Output Display.",
			arguments: arguments{
				debugModeOn: true,
				outputMode:  ctv.MODE_OUTPUT_DISPLAY,
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug False - Output Log.",
			arguments: arguments{
				debugModeOn: false,
				outputMode:  ctv.MODE_OUTPUT_LOG,
			},
			wantError: false,
		},
		{
			name: "Positive Case: Debug True - Output Log.",
			arguments: arguments{
				debugModeOn: true,
				outputMode:  ctv.MODE_OUTPUT_LOG,
			},
			wantError: false,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			tFunctionName, func(t *testing.T) {
				PrintDebugLine(ts.arguments.message, ts.arguments.debugModeOn, ts.arguments.outputMode)
			},
		)
	}
}
