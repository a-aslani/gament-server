package utility

import (
	"Server/app/constants"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//const (
//	privateKeyPath = "keys/app.rsa"
//	publicKeyPath  = "keys/app.rsa.pub"
//)

type TokenClaims struct {
	Key string `json:"key"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type RegisterTokenClaims struct {
	Phone string `json:"phone"`
	jwt.StandardClaims
}

//var (
//	verifyKeyByte, signKeyByte []byte
//	VerifyKey                  *rsa.PublicKey
//	signKey                    *rsa.PrivateKey
//)

// Initialize pub/private keys from path
//func InitKeys() {
//	var err error
//	signKeyByte, err = ioutil.ReadFile(privateKeyPath)
//	if err != nil {
//		log.Fatal("Error reading private key")
//		return
//	}
//
//	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyByte)
//	if err != nil {
//		log.Fatalf("[initKeys]: %s\n", err)
//	}
//
//	verifyKeyByte, err = ioutil.ReadFile(publicKeyPath)
//	if err != nil {
//		log.Fatal("Error reading public key")
//		return
//	}
//
//	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyByte)
//	if err != nil {
//		log.Fatalf("[initKeys]: %s\n", err)
//		panic(err)
//	}
//}

// Create HS512 token
func Token(key string, role string) string {

	//token := jwt.NewWithClaims(jwt.SigningMethodRS256, &TokenClaims{
	//	key,
	//	role,
	//	jwt.StandardClaims{
	//		//ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	//	},
	//})

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &TokenClaims{
		key,
		role,
		jwt.StandardClaims{
			//ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})

	//tokenString, err := token.SignedString(signKey)
	tokenString, err := token.SignedString([]byte(constants.SecretKeyHS512))
	CheckErr(err)
	return tokenString
}

//Create HS256 token
func RegisterToken(phoneKey string) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &RegisterTokenClaims{
		phoneKey,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * constants.ExpiresAtRegisterToken).Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte(constants.SecretKey))
	CheckErr(err)
	return tokenString
}
