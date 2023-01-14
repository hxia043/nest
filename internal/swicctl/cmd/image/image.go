package image

import (
	"fmt"
)

const argsLimitCount int = 1

func getRegistryName(args []string) (string, error) {
	if len(args) != argsLimitCount {
		return "", fmt.Errorf("incorrect number of arguments, the arguments number require %d", argsLimitCount)
	}

	return args[0], nil
}
