package correctionai

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
)

func CorrectionAI(w http.ResponseWriter, r *http.Request) {
	jwt := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", 1))
	if jwt == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ctx := context.Background()
	conf := &firebase.Config{}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Printf("firebase new app: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("firebase auth: %v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	token, err := client.VerifyIDToken(ctx, jwt)
	if err != nil {
		log.Printf("firebase verify ID token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("uid: %v", token.UID)

	/*
		var d struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			fmt.Fprint(w, "Hello, World!")
			return
		}
		if d.Name == "" {
			fmt.Fprint(w, "Hello, World!")
			return
		}*/
	w.WriteHeader(http.StatusOK)
}
