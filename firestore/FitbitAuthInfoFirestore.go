package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

type FitbitAuthInfo struct {
	AccessToken  string
	ExpiresIn 	 int64
	RefreshToken string
	Scope 		 string
	TokenType 	 string
	UserId 		 string
}

func Get(ctx context.Context, client *firestore.Client) (string, error) {
	dsnap, err := client.Collection("roselia").Doc("fitbit").Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	var fitbitAuthInfo FitbitAuthInfo
	err = dsnap.DataTo(&fitbitAuthInfo)

	return fitbitAuthInfo.RefreshToken, err
}
