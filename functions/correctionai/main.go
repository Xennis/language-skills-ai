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
	envAllowOrigin   = "ALLOW_ORIGIN"
	envOpenaiAPIKey  = "OPENAI_API_KEY"
	envOpenaiTesting = "OPENAI_TESTING"

	// See https://beta.openai.com/docs/usage-guidelines/app-review
	// "Please use the end-user ID 'testing' for this activity."
	openaiUIDTesting = "testing"
)

func CorrectionAI(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	allowOrigin, ok := os.LookupEnv(envAllowOrigin)
	if !ok {
		log.Printf("env %s is missing", envAllowOrigin)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", allowOrigin)

	openaiAPIKey, ok := os.LookupEnv(envOpenaiAPIKey)
	if !ok {
		log.Printf("env %s is missing", envOpenaiAPIKey)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	envOpenaiTesting, _ := os.LookupEnv(envOpenaiTesting)

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

	uid := token.UID
	if envOpenaiTesting != "" {
		log.Printf("running uid=%s", openaiUIDTesting)
		uid = openaiUIDTesting
	}

	// skills
	resp, err := skills.Improve(ctx, openaiAPIKey, uid, input.Text)
	if err != nil {
		log.Printf("skills improve: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}
