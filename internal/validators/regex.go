package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// RegexMatches returns a validator that ensures a string matches the given regex pattern
func RegexMatches(pattern *regexp.Regexp, message string) validator.String {
	return stringvalidator.RegexMatches(pattern, message)
}
