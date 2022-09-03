package request

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	firebase_api "ukgm/internal/app/config/firebase"
	"ukgm/internal/app/ukgm/model/entities"
	responses "ukgm/internal/app/ukgm/model/response"

	"github.com/gbrlsnchs/jwt/v2"
)

/*
|--------------------------------------------------------------------------
| pirivat and public key
|--------------------------------------------------------------------------
|
*/

func GetPubPicKey() (pubKey *rsa.PublicKey, privKey *rsa.PrivateKey) {
	var cred *firebase_api.Credential

	privKeyS := firebase_api.GetPrivateKey(firebase_api.CredPath, cred)

	blockPriv, _ := pem.Decode([]byte(privKeyS))
	if blockPriv == nil {
		panic("failed to parse PEM block containing the public key")
	}

	privRes, err := x509.ParsePKCS8PrivateKey(blockPriv.Bytes)
	if err != nil {
		panic("failed to parse RSA encoded privade key: " + err.Error())
	}

	privKey = privRes.(*rsa.PrivateKey)

	var pub *firebase_api.PublicKey

	pubKeyS := firebase_api.GetPublicKey(firebase_api.PubKeyPath, pub)

	block, _ := pem.Decode([]byte(pubKeyS))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}

	// ParsePKIXPublicKey
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Println(err)
		}
	}

	pubRes, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}

	pubKey = pubRes.(*rsa.PublicKey)

	return pubKey, privKey
}

/*
|--------------------------------------------------------------------------
| Middlewares
|--------------------------------------------------------------------------
|
*/
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		token := r.Header.Get("Authorization")
		splitToken := strings.Split(token, "Bearer ")
		if len(token) <= 0 {
			token = r.URL.Query().Get("token")

			if len(token) <= 0 {
				responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("401")), w)
				return
			}
		} else {
			token = splitToken[1]
		}

		pubKey, _ := GetPubPicKey()

		hs256 := jwt.NewRS256(&rsa.PrivateKey{}, pubKey)

		payload, sig, err := jwt.Parse(token)
		if err != nil {
			responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("400-Token-Parse")), w)
			return
		}
		if err = hs256.Verify(payload, sig); err != nil {
			responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("400-Signature")), w)
			return
		}
		var jot Token
		if err = jwt.Unmarshal(payload, &jot); err != nil {
			responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("400-Token-Parse")), w)
			return
		}

		auth := &entities.Auth{
			Auth:  jot.ToAuth(),
			Token: token,
		}

		// Validate fields.
		iatValidator := jwt.IssuedAtValidator(now)
		expValidator := jwt.ExpirationTimeValidator(now)
		audValidator := jwt.AudienceValidator("admin")
		if err = jot.Validate(iatValidator, expValidator, audValidator); err != nil {
			switch err {
			case jwt.ErrIatValidation:
				responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("400-Signature")), w)
				return
			case jwt.ErrExpValidation:
				responses.SendErrorJSON(responses.PrintErrorMessage(errors.New("401-Expired")), w)
				return
			}
		}

		log.Println(auth)

		ctx := context.WithValue(r.Context(), GetAuthCtxKey(), auth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
