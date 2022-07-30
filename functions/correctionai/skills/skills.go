package skills

import (
	"context"
	"fmt"

	"github.com/Xennis/language-skills-ai/functions/correctionai/openai"
)

func Improve(ctx context.Context, openaiAPIKey, user, text string) (string, error) {
	openAIC := openai.NewClient(openaiAPIKey)
	resp, err := openAIC.Completions(ctx, openai.Completions{
		Model:            "text-curie-001",
		Prompt:           []string{"Verbesse dieses in korrektes Deutsch:\n" + text},
		Temperature:      0.0,
		MaxTokens:        60,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		User:             user,
	})
	if err != nil {
		return "", fmt.Errorf("completions: %w", err)
	}
	return resp.Choices[0].Text, nil
}
