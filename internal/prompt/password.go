package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

var stdinReader *bufio.Reader

func getStdinReader() *bufio.Reader {
	if stdinReader == nil {
		stdinReader = bufio.NewReader(os.Stdin)
	}
	return stdinReader
}

func ReadPassword(promptMsg string) (string, error) {
	fmt.Fprint(os.Stderr, promptMsg)

	if term.IsTerminal(int(os.Stdin.Fd())) {
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		return string(password), nil
	}

	reader := getStdinReader()
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Fprintln(os.Stderr)
	return strings.TrimRight(password, "\r\n"), nil
}

func ReadPasswordTwice() (string, error) {
	password, err := ReadPassword("Enter password: ")
	if err != nil {
		return "", err
	}

	if len(password) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}

	confirm, err := ReadPassword("Confirm password: ")
	if err != nil {
		return "", err
	}

	if password != confirm {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}
