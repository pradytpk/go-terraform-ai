package cli

import (
	"flag"
	"log"
	"pradytpk/go-terraform-ai/pkg/terraform"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/walles/env"
)

const version = "0.0.2"

var (
	openAIDeploymentName = flag.String("openai-deployment-name", env.GetOr("OPENAI_DEPLOYMENT_NAME", env.String, "text-davinci-003"), "The deployment name used for the model in OpenAI service.")
	workingDir           = flag.String("working-dir", env.GetOr("WORKING_DIR", env.String, ""), "The path of the project that you want to run")
	execDir              = flag.String("exec-dir", env.GetOr("EXEC_DIR", env.String, ""), "The path of the project that you want to run")
	openAIAPIKey         = flag.String("open-ai-key", env.GetOr("OPENAI_API_KEY", env.String, ""), "The API key for the openai service.This is required")
	requireConfirmation  = flag.Bool("require-confirmation", env.GetOr("REQUIRE_CONFIRMATION", strconv.ParseBool, true), "Whether to require confirmation before executing the command. Defaults to true.")
	temperature          = flag.Float64("temperature", env.GetOr("TEMPERATURE", env.WithBitSize(strconv.ParseFloat, 64), 0.0), "The temperature to use for the model. Range is between 0 and 1. Set closer to 0 if your want output to be more deterministic but less creative. Defaults to 0.0.")
	maxTokens            = flag.Int("max-tokens", env.GetOr("MAX_TOKENS", strconv.Atoi, 0), "The max token will overwrite the max tokens in the max tokens map.")
	azureOpenAIEndpoint  = flag.String("azure-openai-endpoint", env.GetOr("AZURE_OPENAI_ENDPOINT", env.String, ""), "The endpoint for Azure OpenAI service. If provided, Azure OpenAI service will be used instead of OpenAI service.")

	ops terraform.Ops
	err error
)

// InitAndExecute initializes the working directory and execution directory, parses command line flags and executes the root command
//
//	@param workDir
//	@param executionDir
func InitAndExecute(workDir string, executionDir string) {
	flag.Parse()
	if *workingDir == "" {
		workingDir = &workDir
	}
	if *execDir == "" {
		execDir = &executionDir
	}
	if *openAIAPIKey == "" {
		log.Fatal("Please provide an openAI Key")
	}
	if err := RootCmd().Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

// RootCmd starts off the CLI
//
//	@return *cobra.Command
func RootCmd() *cobra.Command {
	ops, err = terraform.NewTerraform(*workingDir, *execDir)
	if err != nil {
		return nil
	}
	cmd := &cobra.Command{
		Use:          "terraform-assistant",
		Version:      version,
		Args:         cobra.MinimumNArgs(1),
		RunE:         runCommand,
		SilenceUsage: true,
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	initCmd := addInit()
	cmd.AddCommand(initCmd)
	return cmd
}
