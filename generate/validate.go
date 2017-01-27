package generate

import (
	"errors"
	"fmt"
	"os"
)

func ValidateArgsLength(genType string, args []string, min int, max int) error {
	length := len(args)

	if length < min {
		if max == min {
			return errors.New(fmt.Sprintf("You must specify %d %s parameter(s) to create", min, genType))
		} else {
			return errors.New(fmt.Sprintf("You must specify at least %d %s parameter(s) to create", min, genType))
		}
	} else if max != -1 && length > max {
		return errors.New(fmt.Sprintf("You may not specify more than %d %s parameter(s) to create at a time", max, genType))
	} else {
		return nil
	}
}

func ValidateNameAvailable(genType string, name string) error {

	if _, err := os.Stat(name); err == nil {
		return errors.New(fmt.Sprintf("%s with name '%s' already exists", genType, name))
	} else {
		return nil
	}

}
