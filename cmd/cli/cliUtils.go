package cli

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	apply     = "Apply"
	dontApply = "Dont Apply"
	reprompt  = "Reprompt"
)

// userActionPrompt
//
//	@return string
//	@return error
func userActionPrompt() (string, error) {
	var (
		result string
		err    error
	)
	if !*requireConfirmation {
		return apply, nil
	}
	items := []string{apply, dontApply}
	lablel := fmt.Sprintf("Would you like to apply this?[%s,%s,%s]", reprompt, items[0], items[1])
	prompt := promptui.SelectWithAdd{
		Label:    lablel,
		Items:    items,
		AddLabel: reprompt,
	}
	_, result, err = prompt.Run()
	if err != nil {
		return dontApply, fmt.Errorf("error to runt the prompt:%w", err)
	}
	return result, nil
}
