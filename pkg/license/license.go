package license

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/yi-cloud/rest-server/pkg/logs"
	"strings"
	"time"
)

type LicenseRef struct {
	Start    string `json:"start,omitempty"`
	End      string `json:"end,omitempty"`
	Clusters int    `json:"clusters,omitempty"`
	SN       string `json:"sn,omitempty"`
	Product  string `json:"product,omitempty"`
}

const (
	IsNone           = "none"
	DecodeError      = "decodeError"
	Expire           = "expire"
	ProductException = "productException"
	Normal           = "normal"
	SNError          = "snError"
	ExceedClusters   = "exceedClusters"
)

var (
	CheckResult = IsNone
	ClusterId   string
	Product     = "ConsoleServer"
)

func decode(license string) (*LicenseRef, error) {
	token, err := jwt.Parse(license, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return RsaPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claim := token.Claims.(jwt.MapClaims)
	ref := &LicenseRef{}
	ref.Start = claim["start"].(string)
	ref.End = claim["end"].(string)
	ref.Clusters = int(claim["clusters"].(float64))
	ref.SN = claim["sn"].(string)
	ref.Product = claim["product"].(string)

	return ref, nil
}

func valid(ref *LicenseRef) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if now < ref.Start || now > ref.End {
		logs.Logger.Warning("License has expired beyond its validity period.")
		CheckResult = Expire
		return
	}

	if ref.Product != Product {
		logs.Logger.Warning("The license is for %s, not valid for %s.", ref.Product, Product)
		CheckResult = ProductException
		return
	}

	validateSN := func(raw_sn string) bool {
		multisn := strings.Split(ref.SN, "|")
		for _, sn := range multisn {
			if sn == raw_sn {
				return true
			}
		}
		return false
	}

	if ClusterId == "" {
		// to get cluster Id
		ClusterId = "12345678"
	}

	if validateSN(ClusterId) {
		CheckResult = Normal
	} else {
		CheckResult = SNError
	}
}

func CheckLicense() {
	license := viper.GetString("server.license")
	if license == "" {
		logs.Logger.Warning("License is none.")
		CheckResult = IsNone
		return
	}

	license = strings.Replace(license, "\n", "", -1)
	go func() {
		for {
			ref, err := decode(license)
			if err != nil {
				logs.Logger.Warning("License is invalid.")
				CheckResult = DecodeError
				return
			}

			valid(ref)
			if CheckResult != Normal {
				return
			}
			time.Sleep(3 * time.Minute)
		}
	}()
}
