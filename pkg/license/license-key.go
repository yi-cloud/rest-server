package license

import (
	"github.com/golang-jwt/jwt/v4"
)

var (
	PublicKey any

	_publicKey string = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEb6XaaYKPxHcypR9JslSoF31eX/h3
qWjZRsi5iM32W90R/t6LcvDs6T5SX9EfE/WdirlXsYSerI9jbXE/mx0dhg==
-----END PUBLIC KEY-----`
)

func init() {
	var err error
	PublicKey, err = jwt.ParseECPublicKeyFromPEM([]byte(_publicKey))
	if err != nil {
		panic(err)
	}
}
