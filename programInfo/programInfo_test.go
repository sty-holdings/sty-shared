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
	"fmt"
	"runtime"
	"testing"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
)

func TestGetFunctionInfo(tPtr *testing.T) {

	var (
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
		tFunctionInfo      FunctionInfo
	)

	type arguments struct {
		level         int
		filenameValue string
	}

	var (
		gotError bool
	)

	tests := []struct {
		name               string
		arguments          arguments
		wantError          bool
		errorMessageFormat string
	}{
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Level 0 Filename Empty",
			arguments: arguments{
				level:         0,
				filenameValue: ctv.VAL_EMPTY,
			},
			errorMessageFormat: EXPECTING_NO_ERROR_FORMAT,
			wantError:          false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Level 1 Filename Empty",
			arguments: arguments{
				level:         1,
				filenameValue: ctv.VAL_EMPTY,
			},
			errorMessageFormat: EXPECTING_NO_ERROR_FORMAT,
			wantError:          false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Level 2 Filename Empty",
			arguments: arguments{
				level:         2,
				filenameValue: ctv.VAL_EMPTY,
			},
			errorMessageFormat: EXPECTING_NO_ERROR_FORMAT,
			wantError:          false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Level 3 Filename Empty",
			arguments: arguments{
				level:         3,
				filenameValue: ctv.VAL_EMPTY,
			},
			errorMessageFormat: EXPECTING_NO_ERROR_FORMAT,
			wantError:          false,
		},
		{
			name: ctv.TEST_POSITIVE_SUCCESS + "Level 4 Filename Empty",
			arguments: arguments{
				level:         4,
				filenameValue: ctv.VAL_EMPTY,
			},
			errorMessageFormat: EXPECTING_NO_ERROR_FORMAT,
			wantError:          true,
		},
	}

	for _, ts := range tests {
		tPtr.Run(
			ts.name, func(t *testing.T) {
				tFunctionInfo = GetFunctionInfo(ts.arguments.level)
				if tFunctionInfo.FileName == ts.arguments.filenameValue ||
					tFunctionInfo.Name == ctv.VAL_EMPTY {
					gotError = true
				} else {
					fmt.Println("FileName: ", tFunctionInfo.FileName)
					gotError = false
				}
				if gotError != ts.wantError {
					tPtr.Errorf(EXPECTING_NO_ERROR_FORMAT, tFunctionName, UNKNOWN)
				}
			},
		)
	}
}

func TestGetProgramInfo(tPtr *testing.T) {

	var (
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
		tProgramInfo       ProgramInfo
	)

	tPtr.Run(
		tFunctionName, func(tPtr *testing.T) {
			tProgramInfo = GetProgramInfo()
			if tProgramInfo.FileName == ctv.VAL_EMPTY ||
				tProgramInfo.NumberCPUs == ctv.VAL_ZERO ||
				tProgramInfo.GoVersion == ctv.VAL_EMPTY {
				tPtr.Errorf(EXPECTING_NO_ERROR_FORMAT, GetFunctionInfo(1).Name, UNKNOWN)
			}
		},
	)
}

func TestGetWorkingDirectory(tPtr *testing.T) {

	var (
		tFunction, _, _, _ = runtime.Caller(0)
		tFunctionName      = runtime.FuncForPC(tFunction).Name()
		tProgramInfo       ProgramInfo
	)

	tPtr.Run(
		tFunctionName, func(tPtr *testing.T) {
			tProgramInfo.WorkingDirectory = GetWorkingDirectory()
			if tProgramInfo.WorkingDirectory == ctv.VAL_EMPTY {
				tPtr.Errorf(EXPECTING_NO_ERROR_FORMAT, GetFunctionInfo(1).Name, UNKNOWN)
			}
		},
	)
}
