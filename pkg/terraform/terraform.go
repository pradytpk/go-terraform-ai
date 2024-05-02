package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// Terraform structure
type Terraform struct {
	WorkingDir string
	ExecDir    string
	Exec       *tfexec.Terraform
}

// NewTerraform creates a new instances of the terraform struct
//
//	@param workingDir
//	@param execDir
//	@return *Terraform
//	@return error
func NewTerraform(workingDir string, execDir string) (*Terraform, error) {
	tf, err := tfexec.NewTerraform(workingDir, execDir)
	if err != nil {
		return nil, fmt.Errorf("Error new terraform:%w", err)
	}
	return &Terraform{
		WorkingDir: workingDir,
		ExecDir:    execDir,
		Exec:       tf,
	}, nil
}
