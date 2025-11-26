package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	PrivateKey    any
	PublicKey     any
	SigningMethod jwt.SigningMethod
)

func InitRsaKey() {
	var err error
	authType := viper.GetString("auth.type")
	if authType == "ecdsa" {
		PrivateKey, err = jwt.ParseECPrivateKeyFromPEM(readKey("auth.privatekey"))
		SigningMethod = jwt.SigningMethodES256
	} else {
		PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(readKey("auth.privatekey"))
		SigningMethod = jwt.SigningMethodRS256
	}
	if err != nil {
		panic(err)
	}

	if authType == "ecdsa" {
		PublicKey, err = jwt.ParseECPublicKeyFromPEM(readKey("auth.publickey"))
	} else {
		PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(readKey("auth.publickey"))
	}
	if err != nil {
		panic(err)
	}

}

func readKey(key string) []byte {
	filename := viper.GetString(key)
	if filename == "" {
		if strings.HasSuffix(key, "privatekey") {
			filename = "/etc/rest-server/auth/private.pem"
		} else {
			filename = "/etc/rest-server/auth/public.pem"
		}
	}
	// read the raw contents of the file
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return data
}
