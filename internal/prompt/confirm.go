package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Confirm(prompt string, defaultYes bool) (bool, error) {
	suffix := " [y/N] "
	if defaultYes {
		suffix = " [Y/n] "
	}

	fmt.Fprint(os.Stderr, prompt+suffix)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))

	if response == "" {
		return defaultYes, nil
	}

	if response == "y" || response == "yes" {
		return true, nil
	}

	if response == "n" || response == "no" {
		return false, nil
	}

	return false, fmt.Errorf("invalid response: %s", response)
}
