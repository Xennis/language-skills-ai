package correctionai

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/Xennis/language-skills-ai/functions/correctionai/openai"
)

func CorrectionAI(w http.ResponseWriter, r *http.Request) {
	jwt := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", 1))
	if jwt == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	openAIKey, ok := os.LookupEnv("OPEN_AI_KEY")
	if !ok {
		log.Printf("env %s is missing", "OPEN_AI_KEY")
		w.WriteHeader(http.StatusInternalServerError)
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
	authC, err := app.Auth(ctx)
	if err != nil {
		log.Printf("firebase auth: %v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	token, err := authC.VerifyIDToken(ctx, jwt)
	if err != nil {
		log.Printf("firebase verify ID token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	openAIC := openai.NewClient(openAIKey)
	resp, err := openAIC.Completion(ctx, openai.Completions{
		Model:            "text-davinci-002",
		Prompt:           []string{"Correct this to standard English:\n\nShe no went to the market."},
		Temperature:      0.0,
		MaxTokens:        60,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		User:             token.UID,
	})
	if err != nil {
		log.Printf("post openai: %v", err)
		return
	}
	fmt.Print(resp.Choices[0].Text)
}
