package prompt

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func WaitForEnter(message string) {
	fmt.Fprintln(os.Stderr, message)

	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Scanln()
	} else {
		reader := getStdinReader()
		reader.ReadString('\n')
	}
}
