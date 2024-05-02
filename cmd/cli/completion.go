package cli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	azureopenai "pradytpk/go-terraform-ai/pkg/gpt3"

	openai "github.com/PullRequestInc/go-gpt3"
	gptEncoder "github.com/samber/go-gpt-3-encoder"
)

// Struct to hold the clients for Azure and OpenAI
type oaiClients struct {
	azureClient  azureopenai.Client
	openAIClient openai.Client
}

var (
	// Map to hold the maximum tokens allowed for different GPT models
	maxTokensMap = map[string]int{
		"code-davinci-002":   8001,
		"text-davinci-003":   4097,
		"gpt-3.5-turbo-0301": 4096,
		"gpt-3.5-turbo":      4096,
		"gpt-35-turbo-0301":  4096, // for azure
		"gpt-4-0314":         8192,
		"gpt-4-32k-0314":     8192,
	}
	// Error for invalid max tokens
	errToken = errors.New("invalid max tokens")
)

const userRole = "user"

// Function to create new OpenAI and Azure clients
func newOAIClients() (oaiClients, error) {
	var (
		oaiClient   openai.Client
		azureClient azureopenai.Client
		err         error
	)
	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		oaiClient = openai.NewClient(*openAIAPIKey)
	} else {
		re := regexp.MustCompile(`^[a-zA-Z0-9]+([_-]?[a-zA-Z0-9]+)*$`)
		if !re.MatchString(*openAIDeploymentName) {
			return oaiClients{}, errors.New("azure openai deployment can only include alphanumeric characters,'_,-', and cant end with '_' or '-'")
		}
		azureClient, err = azureopenai.NewClient(*azureOpenAIEndpoint, *openAIAPIKey, *openAIDeploymentName)
		if err != nil {
			return oaiClients{}, fmt.Errorf("error create azure client:%w", err)
		}
	}
	clients := oaiClients{
		azureClient:  azureClient,
		openAIClient: oaiClient,
	}
	return clients, nil
}

// completion is function that generates completions for given prompt and deployment configuration
//
//	@param ctx
//	@param client
//	@param prompts
//	@param deploymentName
//	@param subcommand
//	@return string
//	@return error
func completion(ctx context.Context, client oaiClients, prompts []string, deploymentName string, subcommand string) (string, error) {
	temp := float32(*temperature)
	maxTokens, err := calculateMaxTokens(prompts, deploymentName)
	if err != nil {
		return "", fmt.Errorf("error prompt string builder:%w", err)
	}
	var prompt strings.Builder
	_, err = fmt.Fprint(&prompt, subcommand)
	if err != nil {
		return "", fmt.Errorf("error prompt string builder:%w", err)
	}
	for _, p := range prompts {
		_, err = fmt.Fprint(&prompt, "%s\n", p)
		if err != nil {
			return "", fmt.Errorf("error range prompt:%w", err)
		}
	}
	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		if isGptTurbo(deploymentName) || isGpt4(deploymentName) {
			res, err := client.openaiGptChatCompletion(ctx, prompt, maxTokens, temp)
			if err != nil {
				return "", fmt.Errorf("error openai gptchart Completion:%w", err)
			}
			return res, nil
		}
		res, err := client.openaiGptCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", fmt.Errorf("error open ai gpt Completion:%w", err)
		}
		return res, nil
	}
	if isGptTurbo35(deploymentName) || isGpt4(deploymentName) {
		resp, err := client.azureGptChatCompletion(ctx, prompt, maxTokens, temp)
		if err != nil {
			return "", fmt.Errorf("error azure GptChat Completion:%w", err)
		}
		return resp, nil
	}
	resp, err := client.azureGptCompletion(ctx, prompt, maxTokens, temp)
	if err != nil {
		return "", fmt.Errorf("error azure Gpt completion: %w", err)
	}

	return resp, nil
}

// calculateMaxTokens is a function that calculates the maximum tokena allowed for a given deployment name
//
//	@param prompts
//	@param deploymentName
//	@return *int
//	@return error
func calculateMaxTokens(prompts []string, deploymentName string) (*int, error) {
	// Get the maximum tokens allowed for the deploymentName from the maxTokensMap
	maxTokensFinal, ok := maxTokensMap[deploymentName]
	if !ok {
		return nil, errors.Wrapf(errToken, "deploymentName %q not found in max tokens map", deploymentName)
	}

	// If a custom maxTokens value is provided, override the value from the map
	if *maxTokens > 0 {
		maxTokensFinal = *maxTokens
	}

	// Create a new gptEncoder
	encoder, err := gptEncoder.NewEncoder()
	if err != nil {
		return nil, fmt.Errorf("error encode gpt: %w", err)
	}

	// Start at 100 since the encoder at times doesn't get it exactly correct
	totalTokens := 100

	// Encode each prompt and calculate the total number of tokens
	for _, prompt := range prompts {
		tokens, err := encoder.Encode(prompt)
		if err != nil {
			return nil, fmt.Errorf("error encode prompt: %w", err)
		}

		totalTokens += len(tokens)
	}

	// Calculate the remaining tokens by subtracting the total tokens from the maximum tokens allowed
	remainingTokens := maxTokensFinal - totalTokens

	return &remainingTokens, nil
}

// isGptTurbo35 is a function that checks deployment names.
//
//	@param deploymentName
//	@return bool
func isGptTurbo35(deploymentName string) bool {
	return deploymentName == "gpt-35-turbo-0301" || deploymentName == "gpt-35-turbo"
}

// isGpt4 is a function that checks deployment names.
//
//	@param deploymentName
//	@return bool
func isGpt4(deploymentName string) bool {
	return deploymentName == "gpt-4-0314" || deploymentName == "gpt-4-32k-0314"
}

// isGptTurbo is a function that checks deployment names.
//
//	@param deploymentName
//	@return bool
func isGptTurbo(deploymentName string) bool {
	return deploymentName == "gpt-3.5-turbo-0301" || deploymentName == "gpt-3.5-turbo"
}
