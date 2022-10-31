package firebase_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"ukgm/internal/app/config/env"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

const (
	CredPath   = "./internal/pkg/token/firebase/ukgm-f77a7-firebase-adminsdk-s2wqn-9226607831.json"
	PubKeyPath = "./internal/pkg/token/other/pubKey.json"

	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
	verifyPasswordURL    = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s"
)

type FirebaseApi struct {
	FirebaseApp *firebase.App
}

func NewFirebaseApi() *FirebaseApi {
	return &FirebaseApi{
		FirebaseApp: SetupFirebase(),
	}
}

func SetupFirebase() *firebase.App {

	serviceAccountKeyFilePath, err := filepath.Abs("./internal/pkg/token/firebase/ukgm-f77a7-firebase-adminsdk-s2wqn-9226607831.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.TODO(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	return app
}

func (f *FirebaseApi) AuthenticationByID(uid string) *auth.UserRecord {
	// Get an auth client from the firebase.App
	client, err := f.FirebaseApp.Auth(context.TODO())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	u, err := client.GetUser(context.TODO(), uid)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	// log.Printf("Successfully fetched user data: %v\n", u)

	return u
}

func (f *FirebaseApi) GenerateToken(uid string) string {
	client, err := f.FirebaseApp.Auth(context.TODO())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.CustomTokenWithClaims(context.TODO(), uid, nil)
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	// log.Printf("Got custom token: %v\n", token)

	return token
}

func (f *FirebaseApi) EncodeToken(token string) (auth *auth.Token, err error) {
	client, err := f.FirebaseApp.Auth(context.TODO())
	if err != nil {
		zap.S().Error(err)
		return nil, fmt.Errorf("error getting auth client: %v", err)
	}

	authTok, err := client.VerifyIDToken(context.TODO(), token)
	if err != nil {
		zap.S().Error(err)
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}

	// log.Printf("Verified ID token: %v\n", authTok)

	return authTok, nil
}

func (f *FirebaseApi) postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		zap.S().Error(err)
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func (f *FirebaseApi) SignInWithCustomToken(token string) (string, error) {
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

	resp, err := f.postRequest(fmt.Sprintf(verifyCustomTokenURL, env.Get("WEB_API_KEY")), req)
	if err != nil {
		zap.S().Error(err)
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		zap.S().Error(err)
		return "", err
	}
	return respBody.IDToken, err
}

func (f *FirebaseApi) SignInWithPassword(email, password string) (string, error) {
	req, err := json.Marshal(map[string]interface{}{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return "", err
	}

	resp, err := f.postRequest(fmt.Sprintf(verifyPasswordURL, env.Get("WEB_API_KEY")), req)
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
