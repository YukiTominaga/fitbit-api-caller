package main

import (
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"github.com/YukiTominaga/fitbit-api-caller/firestore"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RefreshToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	UserId       string `json:"user_id"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/heart-rate", getHeartRateHandler).Methods("GET")
	http.Handle("/", r)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func getHeartRateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "ca-tominaga-flutter"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Read refresh_token
	refreshToken, err := firestore.Get(ctx, client)

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	requestUrl := "https://api.fitbit.com/oauth2/token"
	req, _ := http.NewRequest("POST", requestUrl, strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "Basic MjJESkJQOmZhODlkMTAyYzcxMDZjMmM0NzA2YTJlMzIxNmJmOGYw")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	var token RefreshToken
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&token)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Collection("roselia").Doc("fitbit").Set(ctx, map[string]interface{}{
		"refreshToken": token.RefreshToken,
		"userId": token.UserId,
		"accessToken": token.AccessToken,
		"scope": token.Scope,
		"tokenType": token.TokenType,
		"expiresIn": token.ExpiresIn,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
}