package terraform

import (
	"github.com/pkg/errors"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var errTemplate = errors.New("invalid terraform template")

// CheckTemplate to check the template is valid by parsing it using hclsyntax.ParseConfig
//
//	@param completion
//	@return error
func CheckTemplate(completion string) error {
	template := []byte(completion)
	_, parseDiags := hclsyntax.ParseConfig(template, "", hcl.Pos{Line: 2, Column: 1})
	if len(parseDiags) != 0 {
		return errors.Wrapf(errTemplate, "unexpected valid template but :%s", completion)
	}
	return nil
}
