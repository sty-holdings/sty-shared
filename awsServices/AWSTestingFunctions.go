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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"albert/core/coreFirebase"
	"albert/core/coreFirestore"
	"albert/core/coreHelpers"
	"cloud.google.com/go/firestore"
	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

var (
	// Test Token variables
	testingAccessTokenValid []byte
	testingIdTokenValid     []byte
)

// StartTest - is the AWSServices version which always needs AWSTokens to be current.
func StartTest() (
	myAWS AWSHelper,
	myFireBase coreHelpers.FirebaseFirestoreHelper,
) {

	var (
		err       error
		errorInfo pi.ErrorInfo
	)

	myAWS, errorInfo = NewAWSSession(ctv.TEST_AWS_INFORMATION_FQN)

	if errorInfo.Error == nil {
		myFireBase.CredentialsLocation = ctv.TEST_FIREBASE_CREDENTIALS
		if myFireBase.AppPtr, myFireBase.AuthPtr, errorInfo = coreFirebase.GetFirebaseAppAuthConnection(myFireBase.CredentialsLocation); errorInfo.Error == nil {
			myFireBase.FirestoreClientPtr, errorInfo = coreFirestore.GetFirestoreClientConnection(myFireBase.AppPtr)
		}
	}

	BuildTestUser(myFireBase)
	if err = loadTestingTokens(myAWS, myFireBase.FirestoreClientPtr); err != nil {
		os.Exit(1)
	}

	return
}

func StopTest(myFireBase coreHelpers.FirebaseFirestoreHelper) {

	_ = coreFirestore.RemoveDocumentById(myFireBase.FirestoreClientPtr, ctv.DATASTORE_USERS, ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE)

}

func BuildTestUser(myFireBase coreHelpers.FirebaseFirestoreHelper) {

	testUser := map[any]interface{}{
		ctv.FN_REQUESTOR_ID:     ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID,
		ctv.FN_FIRST_NAME:       ctv.TEST_USER_FIRST_NAME,
		ctv.FN_LAST_NAME:        ctv.TEST_USER_LAST_NAME,
		ctv.FN_EMAIL:            ctv.TEST_USER_EMAIL,
		ctv.FN_AREA_CODE:        ctv.TEST_USER_AREA_CODE,
		ctv.FN_PHONE_NUMBER:     ctv.TEST_USER_PHONE_NUMBER,
		ctv.FN_USERNAME:         ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE,
		ctv.FN_CREATE_TIMESTAMP: time.Now(),
	}

	_ = coreFirestore.SetDocument(myFireBase.FirestoreClientPtr, ctv.DATASTORE_USERS, ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID, testUser)

}

func GetValidTestingAccessToken() string {

	return string(testingAccessTokenValid)
}

func GetValidTestingIdToken() string {

	return string(testingIdTokenValid)
}

func loadTestingTokens(
	myAWS AWSHelper,
	firestoreClientPtr *firestore.Client,
) error {

	//goland:noinspection ALL
	const (
		AWS_TEST_LOGIN_UI = "https://savup-test.auth.us-west-2.amazoncognito.com/oauth2/authorize?client_id=5g9kg0c3mnpd5mlaf28ib2m6d9&response_type=token&scope=aws.cognito.signin.user.admin+email+openid&redirect_uri=https%3A%2F%2Fjwt.io"
		AWS_TOKENS_FILE   = "/Users/syacko/workspace/styh-dev/src/albert/core/testTokens/AWSTokens"
	)

	var (
		err        error
		tAWSTokens []byte
		tNameValue []string
		tTokens    []string
		tValid     bool
	)

	// Load raw token files
	if tAWSTokens, err = os.ReadFile(ctv.TEST_AWS_TEST_TOKEN_FQN); err != nil {
		err = errors.New(fmt.Sprintf("Not able to read the AWS Token file: %v.", ctv.TEST_AWS_TEST_TOKEN_FQN))
		log.Println(err.Error())
	} else {
		tAWSTokens = tAWSTokens[16:]
		tTokens = strings.Split(string(tAWSTokens), "&")
		for _, token := range tTokens {
			tNameValue = strings.Split(token, "_token=")
			if tNameValue[0] == ctv.TOKEN_TYPE_ACCESS {
				if err = os.WriteFile(ctv.TEST_AWS_RAW_ACCESS_TOKEN_FQN, []byte(tNameValue[1]), 0666); err != nil {
					err = errors.New(fmt.Sprintf("Not able to write the AWS Access Token file: %v.", ctv.TEST_AWS_RAW_ACCESS_TOKEN_FQN))
					log.Println(err.Error())
				}
			}
			if tNameValue[0] == ctv.TOKEN_TYPE_ID {
				if err = os.WriteFile(ctv.TEST_AWS_RAW_ID_TOKEN_FQN, []byte(tNameValue[1]), 0666); err != nil {
					err = errors.New(fmt.Sprintf("Not able to write the AWS Id Token file: %v.", ctv.TEST_AWS_RAW_ID_TOKEN_FQN))
					log.Println(err.Error())
				}
			}
		}
	}

	// Checking that test ACCESS token is valid
	if testingAccessTokenValid, err = os.ReadFile(ctv.TEST_AWS_RAW_ACCESS_TOKEN_FQN); err == nil {
		if tValid, _ = myAWS.ValidAWSJWT(firestoreClientPtr, ctv.TOKEN_TYPE_ACCESS, string(testingAccessTokenValid)); tValid == false {
			err = errors.New("ACCESS Token loading failed")
			fmt.Println(ctv.EMPTY)
			fmt.Printf("NOTE: Make sure to create user %v in the USERS datastorel.\n", ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE)
			fmt.Println(ctv.EMPTY)
			fmt.Printf("NOTE: The username must be %v or tests will fail.\n", ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE)
			fmt.Printf(
				"The access token is not valid! You need to create one using \n%v and paste the URL into the \n%v file.\n",
				AWS_TEST_LOGIN_UI,
				AWS_TOKENS_FILE,
			)
		}
	}

	// Checking that test ID token is valid
	if testingIdTokenValid, err = os.ReadFile(ctv.TEST_AWS_RAW_ID_TOKEN_FQN); err == nil {
		if tValid, _ = myAWS.ValidAWSJWT(firestoreClientPtr, ctv.TOKEN_TYPE_ID, string(testingIdTokenValid)); tValid == false {
			err = errors.New("ID Token loading failed")
			fmt.Println(ctv.EMPTY)
			fmt.Printf("NOTE: Make sure to create user %v in the USERS datastorel.\n", ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE)
			fmt.Println(ctv.EMPTY)
			fmt.Printf("NOTE: The username must be %v or tests will fail.\n", ctv.TEST_USERNAME_SAVUP_TEST_DO_NOT_DELETE)
			fmt.Printf(
				"The Id token is not valid! You need to create one using (%v) and paste the URL into the %v file.\n",
				AWS_TEST_LOGIN_UI,
				AWS_TOKENS_FILE,
			)
		}
	}

	return err
}

func RemoveTestUser(myFireBase coreHelpers.FirebaseFirestoreHelper) {
	_ = coreFirestore.RemoveDocumentById(myFireBase.FirestoreClientPtr, ctv.DATASTORE_USERS, ctv.TEST_USERNAME_SAVUP_REQUESTOR_ID)
}

func SetValidTestingAccessToken(accessToken []byte) {

	testingAccessTokenValid = accessToken
}

func SetValidTestingIdToken(accessToken []byte) {

	testingIdTokenValid = accessToken
}

func getToken(tokenType, token string) string {

	switch strings.ToUpper(token) {
	case ctv.INVALID:
		return ctv.TEST_STRING
	case ctv.VALID:
		if tokenType == ctv.TOKEN_TYPE_ACCESS {
			return GetValidTestingAccessToken()
		} else {
			return GetValidTestingIdToken()
		}
	case ctv.MISSING:
		return ctv.EMPTY
	}

	return ctv.EMPTY
}
