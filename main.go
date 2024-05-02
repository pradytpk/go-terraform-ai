package main

import (
	"log"
	"pradytpk/go-terraform-ai/cmd/cli"
	"pradytpk/go-terraform-ai/pkg/utils"
)

// main initialize the working directory and terraform executable directory, then calls the InitAndExecute function to start the program
func main() {
	workingDir, err := utils.CurrenDir()
	if err != nil {
		log.Fatalf("failed to get the current dir:%s\n")
	}
	execDir, err := utils.TerraformPath()
	if err != nil {
		log.Fatalf("failed to get the exec dir:%s\n")

	}
	cli.InitAndExecute(workingDir, execDir)
}
