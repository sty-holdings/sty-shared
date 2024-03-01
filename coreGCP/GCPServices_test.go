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
)

var (
//goland:noinspection ALL
)

func TestCreateStorageClient(tPtr *testing.T) {

	// var (
	// 	errorInfo          pi.ErrorInfo
	// 	tFunction, _, _, _ = runtime.Caller(0)
	// 	tFunctionName      = runtime.FuncForPC(tFunction).Name()
	// 	tClient            *storage.Client
	// )
	//
	// tPtr.Run(
	// 	tFunctionName, func(t *testing.T) {
	// 		if tClient, errorInfo = CreateStorageClient(ctv.TEST_GCP_CREDENTIALS, true); tClient == nil || errorInfo.Error != nil {
	// 			tPtr.Errorf("%v Failed: Was expecting file data but got %v.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tClient, errorInfo = CreateStorageClient(ctv.TEST_GCP_CREDENTIALS_INVALID, true); tClient != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got data.", tFunctionName, pi.ERROR)
	// 		}
	// 	},
	// )
}

func TestGetBucket(tPtr *testing.T) {

	// var (
	// 	errorInfo          pi.ErrorInfo
	// 	tBucketPtr         *storage.BucketHandle
	// 	tClient            *storage.Client
	// 	tFunction, _, _, _ = runtime.Caller(0)
	// 	tFunctionName      = runtime.FuncForPC(tFunction).Name()
	// )

	// tClient, _ = CreateStorageClient(ctv.TEST_GCP_CREDENTIALS, true)
	//
	// tPtr.Run(
	// 	tFunctionName, func(t *testing.T) {
	// 		if tBucketPtr, errorInfo = getBucket(tClient, "savup-private"); tBucketPtr == nil || errorInfo.Error != nil {
	// 			tPtr.Errorf("%v Failed: Was expecting a pointer but got %v.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tBucketPtr, errorInfo = getBucket(tClient, ""); tBucketPtr != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got a pointer.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tBucketPtr, errorInfo = getBucket(nil, "savup-private"); tBucketPtr != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got a pointer.", tFunctionName, pi.ERROR)
	// 		}
	// 	},
	// )
}

func TestGetGCPKey(tPtr *testing.T) {

	// var (
	// 	tFunction, _, _, _ = runtime.Caller(0)
	// 	tFunctionName      = runtime.FuncForPC(tFunction).Name()
	// 	tGCPCredentials    []byte
	// )
	//
	// tPtr.Run(
	// 	tFunctionName, func(t *testing.T) {
	// 		if tGCPCredentials = getGCPKey(ctv.TEST_GCP_CREDENTIALS, true); tGCPCredentials == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting file data but got %v.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tGCPCredentials = getGCPKey(ctv.TEST_GCP_CREDENTIALS_INVALID, true); tGCPCredentials != nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got data.", tFunctionName, pi.ERROR)
	// 		}
	// 	},
	// )
}

func TestListObjectsInBucket(tPtr *testing.T) {

	// var (
	// 	errorInfo          pi.ErrorInfo
	// 	tBucketList        []string
	// 	tClient            *storage.Client
	// 	tFunction, _, _, _ = runtime.Caller(0)
	// 	tFunctionName      = runtime.FuncForPC(tFunction).Name()
	// )
	//
	// tClient, _ = CreateStorageClient(ctv.TEST_GCP_CREDENTIALS, true)
	//
	// tPtr.Run(
	// 	tFunctionName, func(t *testing.T) {
	// 		if tBucketList, errorInfo = ListObjectsInBucket(tClient, "savup-private"); tBucketList == nil || errorInfo.Error != nil {
	// 			tPtr.Errorf("%v Failed: Was expecting file data but got %v.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tBucketList, errorInfo = ListObjectsInBucket(tClient, ""); tBucketList != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got data.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tBucketList, errorInfo = ListObjectsInBucket(nil, "savup-private"); tBucketList != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got data.", tFunctionName, pi.ERROR)
	// 		}
	// 	},
	// )
}

func TestReadBucketObject(tPtr *testing.T) {

	// var (
	// 	errorInfo          pi.ErrorInfo
	// 	tClient            *storage.Client
	// 	tContents          []byte
	// 	tFunction, _, _, _ = runtime.Caller(0)
	// 	tFunctionName      = runtime.FuncForPC(tFunction).Name()
	// )
	//
	// tClient, _ = CreateStorageClient(ctv.TEST_GCP_CREDENTIALS, true)
	//
	// tPtr.Run(
	// 	tFunctionName, func(t *testing.T) {
	// 		if tContents, errorInfo = ReadBucketObject(
	// 			tClient,
	// 			"savup-private",
	// 			"templates/promissoryNote/California/SavUp-Promissory-Note.html",
	// 		); tContents == nil || errorInfo.Error != nil {
	// 			tPtr.Errorf("%v Failed: Was expecting file data but got %v.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tContents, errorInfo = ReadBucketObject(
	// 			tClient,
	// 			"",
	// 			"templates/promissoryNote/California/SavUp-Promissory-Note.html",
	// 		); tContents != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got data.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tContents, errorInfo = ReadBucketObject(tClient, "savup-private", ""); tContents != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got data.", tFunctionName, pi.ERROR)
	// 		}
	// 		if tContents, errorInfo = ReadBucketObject(
	// 			nil,
	// 			"savup-private",
	// 			"templates/promissoryNote/California/SavUp-Promissory-Note.html",
	// 		); tContents != nil || errorInfo.Error == nil {
	// 			tPtr.Errorf("%v Failed: Was expecting an error %v but got data.", tFunctionName, pi.ERROR)
	// 		}
	// 	},
	// )
}
