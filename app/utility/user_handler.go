package utility

import (
	"Server/app/constants"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

//Generate hash password
func GenerateHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	CheckErr(err)
	return string(hash)
}

//Verify hash password
func VerifyPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//Check active code is timeout or no
func CheckActiveTime(time, now int64) bool {
	diff := now - time
	if diff <= constants.CodeActiveTime {
		return  true
	}
	return  false
}

//Generate random number
func GenerateRandomCode() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(99999-10000) + 10000
}
