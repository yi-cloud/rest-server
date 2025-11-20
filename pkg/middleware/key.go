package middleware

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	RsaPrivateKey *rsa.PrivateKey
	RsaPublicKey  *rsa.PublicKey
)

func InitRsaKey() {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(readKey("auth.rsaprivatekey"))
	if err != nil {
		panic(err)
	}
	RsaPrivateKey = privateKey

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(readKey("auth.rsapublickey"))
	if err != nil {
		panic(err)
	}
	RsaPublicKey = publicKey

}

func readKey(key string) []byte {
	filename := viper.GetString(key)
	if filename == "" {
		if strings.HasSuffix(key, "rsaprivatekey") {
			filename = "/etc/rest-server/auth/pri.key"
		} else {
			filename = "/etc/rest-server/auth/pub.key"
		}
	}
	// read the raw contents of the file
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return data
}
