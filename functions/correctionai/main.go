package correctionai

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	firebase "firebase.google.com/go"
	"github.com/Xennis/language-skills-ai/functions/correctionai/skills"
)

const (
	envOpenaiAPIKey = "OPENAI_API_KEY"
)

func CorrectionAI(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	openaiAPIKey, ok := os.LookupEnv(envOpenaiAPIKey)
	if !ok {
		log.Printf("env %s is missing", envOpenaiAPIKey)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwt := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	if jwt == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// input
	var input struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if input.Text == "" || utf8.RuneCountInString(input.Text) > 200 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify user
	app, err := firebase.NewApp(ctx, &firebase.Config{})
	if err != nil {
		log.Printf("firebase new app: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	authC, err := app.Auth(ctx)
	if err != nil {
		log.Printf("firebase auth: %v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	token, err := authC.VerifyIDToken(ctx, jwt)
	if err != nil {
		log.Printf("firebase verify id token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// skills
	resp, err := skills.Improve(ctx, openaiAPIKey, token.UID, input.Text)
	if err != nil {
		log.Printf("skills improve: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}
