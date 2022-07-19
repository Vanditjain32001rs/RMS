package utilities

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	passHash, hashErr := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(passHash), hashErr
}

func EncodeToJson(i ...interface{}) ([]byte, error) {
	jsonData, jsonErr := json.Marshal(i)

	return jsonData, jsonErr
}

func Contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
