package marvin

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// ValidateHexString - returns true iff all chars in string are valid hex chars
func ValidateHexString(hex string) bool {
	for _, c := range hex {
		if !strings.Contains("0123456789abcdefABCDEF", string(c)) {
			return false
		}
	}
	return true
}

// ValidateHashType - returns true iff hashType argument is a supported type
func ValidateHashType(hashType string) bool {
	switch hashType {
	case "md5", "md4", "sha1", "sha256", "sha384", "sha512", "ntlm":
		return true
	}
	return false
}

// RemoteHashLookup - Looks up given hash in the md5decrypt online service
func RemoteHashLookup(hash, hashType, email, code string) (string, error) {

	if !ValidateHashType(hashType) {
		return "", errors.New("Invalid Hash Type")
	}

	if !ValidateHexString(hash) {
		return "", errors.New("Invalid Hash String")
	}

	if len(code) < 1 || len(email) < 1 {
		return "", errors.New("md5decrypt API not configured")
	}

	queryString := "https://md5decrypt.net/en/Api/api.php?hash=" + hash +
		"&hash_type=" + hashType + "&email=" + email + "&code=" + code

	resp, err := http.Get(queryString)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body), err
}
