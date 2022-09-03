package firebase_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

const (
	CredPath   = "./internal/pkg/token/firebase/ukgm-f77a7-firebase-adminsdk-s2wqn-9226607831.json"
	PubKeyPath = "./internal/pkg/token/other/pubKey.json"

	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
	verifyPasswordURL    = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s"
)

type Credential struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
}

type PublicKey struct {
	PubKey string `json:"ed6bc9da1af236f8ea56a5d6223401ff4e087f11"`
}

func GetPrivateKey(path string, c *Credential) *Credential {
	configFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&c); err != nil {
		log.Println(err)
	}
	return c
}

func GetPublicKey(path string, p *PublicKey) string {
	configFile, err := os.Open(path)
	if err != nil {
		return err.Error()
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&p); err != nil {
		return err.Error()
	}
	return p.PubKey
}

func SetupFirebase() *firebase.App {

	serviceAccountKeyFilePath, err := filepath.Abs("./internal/pkg/token/firebase/ukgm-f77a7-firebase-adminsdk-s2wqn-9226607831.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	return app
}

func AuthenticationID(app *firebase.App, uid string) *auth.UserRecord {
	// Get an auth client from the firebase.App
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	u, err := client.GetUser(context.Background(), uid)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	log.Printf("Successfully fetched user data: %v\n", u)

	return u
}

func GenerateToken(app *firebase.App, uid string) string {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.CustomTokenWithClaims(context.Background(), uid, nil)
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	log.Printf("Got custom token: %v\n", token)

	return token
}

func EncodeToken(app *firebase.App, token string) *auth.Token {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	authTok, err := client.VerifyIDToken(context.Background(), token)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", authTok)

	return authTok
}

func SignInWithCustomToken(token string, privateKeyID string) (string, error) {
	req, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", err
	}

	// apiKey, err := internal.APIKey()
	// if err != nil {
	// 	return "", err
	// }
	resp, err := postRequest(fmt.Sprintf(verifyCustomTokenURL, privateKeyID), req)
	if err != nil {
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return "", err
	}
	return respBody.IDToken, err
}

func postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
