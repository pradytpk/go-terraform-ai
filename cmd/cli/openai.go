package cli

import (
	"context"
	"fmt"
	"pradytpk/go-terraform-ai/pkg/utils"
	"strings"

	azureopenai "pradytpk/go-terraform-ai/pkg/gpt3"

	openai "github.com/PullRequestInc/go-gpt3"
	"github.com/pkg/errors"
)

var errResp = errors.New("invalid response")

// openaiGptCompletion generates a GPT-3 completion using the OpenAI API.
//
//	@receiver c
//	@param ctx
//	@param prompt
//	@param maxTokens
//	@param temp
//	@return string
//	@return error
func (c *oaiClients) openaiGptCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	resp, err := c.openAIClient.CompletionWithEngine(ctx, *openAIDeploymentName, openai.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           utils.ToPtr(1),
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error openai completion:%w", err)
	}
	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	// Return the generated completion text
	return resp.Choices[0].Text, nil
}

// openaiGptChatCompletion generates a GPT-3 completion using the OpenAI API.
//
//	@receiver c
//	@param ctx
//	@param prompt
//	@param maxTokens
//	@param temp
//	@return string
//	@return error
func (c *oaiClients) openaiGptChatCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	resp, err := c.openAIClient.ChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: *openAIDeploymentName,
		Messages: []openai.ChatCompletionRequestMessage{
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           1,
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error openai gpt completion:%w", err)
	}
	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Message.Content, nil

}

// azureGptCompletion generates a GPT-3 completion using the Azure API.
//
//	@receiver c
//	@param ctx
//	@param prompt
//	@param maxTokens
//	@param temp
//	@return string
//	@return error
func (c *oaiClients) azureGptCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	resp, err := c.azureClient.Completion(ctx, azureopenai.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           utils.ToPtr(1),
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error azure completion: %w", err)
	}

	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Text, nil
}

// azureGptChatCompletion generates a GPT-3 completion using the Azure API.
//
//	@receiver c
//	@param ctx
//	@param prompt
//	@param maxTokens
//	@param temp
//	@return string
//	@return error
func (c *oaiClients) azureGptChatCompletion(ctx context.Context, prompt strings.Builder, maxTokens *int, temp float32) (string, error) {
	resp, err := c.azureClient.ChatCompletion(ctx, azureopenai.ChatCompletionRequest{
		Model: *openAIDeploymentName,
		Messages: []azureopenai.ChatCompletionRequestMessage{
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           1,
		Temperature: &temp,
	})
	if err != nil {
		return "", fmt.Errorf("error azure chatgpt completion: %w", err)
	}

	if len(resp.Choices) != 1 {
		return "", errors.Wrapf(errResp, "expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Message.Content, nil
}
