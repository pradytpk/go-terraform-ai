package terraform

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// Init initializes the terraform instances
//
//	@receiver ter
//	@return error
func (ter *Terraform) Init() error {
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Start()

	err := ter.Exec.Init(context.Background())
	if err != nil {
		spin.Stop()
		return fmt.Errorf("error running init:%w", err)
	}
	spin.Stop()
	return nil
}

// Apply applies the Terraform configuration
//
//	@receiver ter
//	@return error
func (ter *Terraform) Apply() error {
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Start()

	err := ter.Exec.Apply(context.Background())
	if err != nil {
		spin.Stop()
		return fmt.Errorf("error running apply:%w", err)
	}
	spin.Stop()
	return nil
}
