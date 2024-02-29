// Package coreHelpersValidators
/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS

    {Other catagories of restrictions}
    * {List of restrictions for the catagory

NOTES:
    {Enter any additional notes that you believe will help the next developer.}

COPYRIGHT:
	Copyright 2022
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.

*/
package coreHelpersValidators

//goland:noinspection GoSnakeCaseUsage
const (
	TEST_BASE64_STRING          = "VEhpcyBpcyBhIHRlc3Qgc3RyaW5nIDEyMzQxMzQ1MjM1Nl4lKl4mJSYqKCVeKg=="
	TEST_FILE_NAME              = "test_file.txt"
	TEST_DIRECTORY              = "/tmp"
	TEST_DIRECTORY_ENDING_SLASH = "/tmp/"
	TEST_DIRECTORY_NON_ROOT     = "shared-services"
	TEST_STRING                 = "THis is a test string 123413452356^%*^&%&*(%^*"
)

var (
	TEST_BYTE_ARRAY = []byte(TEST_STRING)
)

//goland:noinspection GoSnakeCaseUsage
const (
	TEST_DOMAIN               = "savup.com"
	TEST_FILE_EXISTS_FILENAME = "file_exists.txt"
	TEST_FILE_UNREADABLE      = "unreadable_file.txt"
	TEST_INVALID_DOMAIN       = "tmp"
)

var (
	testValidJson = []byte("{\"name\": \"Test Name\"}")
)
