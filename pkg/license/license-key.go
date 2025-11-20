package license

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
)

var (
	RsaPublicKey *rsa.PublicKey

	PublicKey string = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC4B5BfRrsXH2OXq/nILkMMYerU
EoCcz/suR2GIPfSBU6dRDzdTrBQ4BbR5kojJrgKzlziLrqgLM8mlL1ukwc2roV5I
wWbisJD0C5Jqw2LJj66Qs+0iUJsEe3lz/8FosnS28Vj4aIW7Mne2lZaMSygDosME
oWS9wWRmC/BRtrg20QIDAQAB
-----END PUBLIC KEY-----`
)

func init() {
	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(PublicKey))
	if err != nil {
		panic(err)
	}
	RsaPublicKey = rsaPublicKey

}
