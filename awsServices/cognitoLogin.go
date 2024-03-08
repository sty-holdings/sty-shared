package sty_shared

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	ctv "github.com/sty-holdings/constant-type-vars-go/v2024"
	pi "github.com/sty-holdings/sty-shared/v2024/programInfo"
)

// NewCognitoLogin - creates a CognitoLogin object. If you have a clientSecret, we use a pointer
// so there is only one place in memory (Security).
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func NewCognitoLogin(
	username, password, userPoolId, clientId string,
	clientSecret *string,
) (
	*CognitoLogin,
	error,
) {
	c := &CognitoLogin{
		username:     username,
		password:     password,
		userPoolId:   userPoolId,
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	if !strings.Contains(userPoolId, "_") {
		return nil, fmt.Errorf("invalid Cognito User Pool ID (%s), must be in format: '<region>_<pool name>'", userPoolId)
	}
	c.userPoolName = strings.Split(userPoolId, "_")[1]

	c.bigN = hexToBig(nHex)
	c.g = hexToBig(gHex)
	c.k = hexToBig(hexHash("00" + nHex + "0" + gHex))
	c.a = c.generateRandomSmallA()
	c.bigA = c.calculateA()

	return c, nil
}

// GetUsername - returns the configured Cognito user username
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) GetUsername() string {
	return csrp.username
}

// GetClientId - returns the configured Cognito Cient ID
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) GetClientId() string {
	return csrp.clientId
}

// GetUserPoolId - returns the configured Cognito User Pool ID
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) GetUserPoolId() string {
	return csrp.userPoolId
}

// GetUserPoolName - returns the configured Cognito User Pool Name
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) GetUserPoolName() string {
	return csrp.userPoolName
}

// GetAuthParams - returns the Auth Parameter map of values. If getSecret = true, then the
// client secret will be returned. If the client secret is missing, then authParams will be empty.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) GetAuthParams() map[string]string {
	params := map[string]string{
		"USERNAME": csrp.username,
		"PASSWORD": csrp.password,
		"SRP_A":    bigToHex(csrp.bigA),
	}

	if secret, err := csrp.GetSecretHash(csrp.username); err == nil {
		params["SECRET_HASH"] = secret
	}

	return params
}

// GetSecretHash returns the secret hash string required to make certain
// Cognito Identity Provider API calls (if client is configured with a secret)
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) GetSecretHash(username string) (
	string,
	error,
) {
	if csrp.clientSecret == nil {
		return "", fmt.Errorf("unable to create secret hash as client secret has not been configured")
	}

	var (
		msg = username + csrp.clientId
		key = []byte(*csrp.clientSecret)
		h   = hmac.New(sha256.New, key)
	)

	h.Write([]byte(msg))

	sh := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return sh, nil
}

// PasswordVerifierChallenge - returns the ChallengeResponses map to be used
// inside the cognitoidentityprovider.RespondToAuthChallengeInput object which
// fulfils the PASSWORD_VERIFIER Cognito challenge
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) PasswordVerifierChallenge(
	challengeParms map[string]string,
	ts time.Time,
) (
	map[string]string,
	error,
) {
	var (
		internalUsername = challengeParms["USERNAME"]
		userId           = challengeParms["USER_ID_FOR_SRP"]
		saltHex          = challengeParms["SALT"]
		srpBHex          = challengeParms["SRP_B"]
		secretBlockB64   = challengeParms["SECRET_BLOCK"]

		timestamp = ts.In(time.UTC).Format("Mon Jan 2 03:04:05 MST 2006")
		hkdf      = csrp.getPasswordAuthenticationKey(userId, csrp.password, hexToBig(srpBHex), hexToBig(saltHex))
	)

	secretBlockBytes, err := base64.StdEncoding.DecodeString(secretBlockB64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode challenge parameter 'SECRET_BLOCK', %s", err.Error())
	}

	msg := csrp.userPoolName + userId + string(secretBlockBytes) + timestamp
	hmacObj := hmac.New(sha256.New, hkdf)
	hmacObj.Write([]byte(msg))
	signature := base64.StdEncoding.EncodeToString(hmacObj.Sum(nil))

	response := map[string]string{
		"TIMESTAMP":                   timestamp,
		"USERNAME":                    internalUsername,
		"PASSWORD_CLAIM_SECRET_BLOCK": secretBlockB64,
		"PASSWORD_CLAIM_SIGNATURE":    signature,
	}
	if secret, err := csrp.GetSecretHash(internalUsername); err == nil {
		response["SECRET_HASH"] = secret
	}

	return response, nil
}

// generateRandomSmallA - creates the smallA needed for cryptographic security.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) generateRandomSmallA() *big.Int {

	return big.NewInt(0).Mod(getRandom(128), csrp.bigN)
}

// calculateA - creates value A using cryptographic security algorithm. The function '
// then checks if A is divisible by the modulus (a weakness) and generates an error message
// if so, otherwise it returns the calculated A.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) calculateA() *big.Int {

	bigA := big.NewInt(0).Exp(csrp.g, csrp.a, csrp.bigN)
	if big.NewInt(0).Mod(bigA, csrp.bigN).Cmp(big.NewInt(0)) == 0 {
		pi.PrintErrorInfo(pi.NewErrorInfo(pi.ErrNotDivisibleN, ctv.VAL_EMPTY))
	}

	return bigA
}

// getPasswordAuthenticationKey - calculates the auth key using HKDF
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func (csrp *CognitoLogin) getPasswordAuthenticationKey(
	username, password string,
	bigB, salt *big.Int,
) []byte {

	var (
		userPass     = fmt.Sprintf("%s%s:%s", csrp.userPoolName, username, password)
		userPassHash = hashSha256([]byte(userPass))

		uVal      = calculateU(csrp.bigA, bigB)
		xVal      = hexToBig(hexHash(padHex(salt.Text(16)) + userPassHash))
		gModPowXN = big.NewInt(0).Exp(csrp.g, xVal, csrp.bigN)
		intVal1   = big.NewInt(0).Sub(bigB, big.NewInt(0).Mul(csrp.k, gModPowXN))
		intVal2   = big.NewInt(0).Add(csrp.a, big.NewInt(0).Mul(uVal, xVal))
		sVal      = big.NewInt(0).Exp(intVal1, intVal2, csrp.bigN)
	)

	return computeHKDF(padHex(sVal.Text(16)), padHex(bigToHex(uVal)))
}

// hashShe256 - coverts a byte array to sha256 encoded string
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func hashSha256(buf []byte) string {
	a := sha256.New()
	a.Write(buf)

	return hex.EncodeToString(a.Sum(nil))
}

// hexHash - decodes the hex string to a byte array
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func hexHash(hexStr string) string {

	var (
		errorInfo   pi.ErrorInfo
		tHexStrByte []byte
	)

	if tHexStrByte, errorInfo.Error = hex.DecodeString(hexStr); errorInfo.Error != nil {
		pi.PrintErrorInfo(pi.NewErrorInfo(pi.ErrDecodeStringFailed, fmt.Sprintf("%v%v", ctv.TXT_HEX_STRING, hexStr)))
	}

	return hashSha256(tHexStrByte)
}

// hexToBig - converts a hex string to a big Int pointer
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func hexToBig(hexStr string) *big.Int {

	var (
		i  *big.Int
		ok bool
	)

	if i, ok = big.NewInt(0).SetString(hexStr, 16); ok == false {
		pi.PrintErrorInfo(pi.NewErrorInfo(pi.ErrSetStringFailed, fmt.Sprintf("%v%v", ctv.TXT_HEX_STRING, hexStr)))
	}

	return i
}

// bigToHex - converts a big int pointer to a hex string
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func bigToHex(val *big.Int) string {
	return val.Text(16)
}

// getRandom - creates a random number using rand.Read
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func getRandom(n int) *big.Int {

	b := make([]byte, n)
	rand.Read(b)

	return hexToBig(hex.EncodeToString(b))
}

// padHex - pads a big.Int with leading zeros (replace if needed)
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func padHex(hexStr string) string {

	if len(hexStr)%2 == 1 {
		return fmt.Sprintf("0%s", hexStr)
	}

	if strings.Contains("89ABCDEFabcdef", string(hexStr[0])) {
		return fmt.Sprintf("00%s", hexStr)
	}

	return hexStr
}

// computeHKDF - uses the standard HKDF algorithm
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func computeHKDF(ikm, salt string) []byte {

	ikmb, _ := hex.DecodeString(ikm)
	saltb, _ := hex.DecodeString(salt)

	extractor := hmac.New(sha256.New, saltb)
	extractor.Write(ikmb)
	prk := extractor.Sum(nil)
	infoBitsUpdate := append([]byte(infoBits), byte(1))
	extractor = hmac.New(sha256.New, prk)
	extractor.Write(infoBitsUpdate)
	hmacHash := extractor.Sum(nil)

	return hmacHash[:16]
}

// calculateU - creates the hash of A and B
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func calculateU(bigA, bigB *big.Int) *big.Int {
	return hexToBig(hexHash(padHex(bigA.Text(16)) + padHex(bigB.Text(16))))
}
