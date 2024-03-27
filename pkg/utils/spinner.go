package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

const spinnerTime = 100 * time.Millisecond

/* StartSpinner starts a spinner with a message and then replaces it with a done message
 * when the spinner is stopped.
 */
func StartSpinner(msg string, finalMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[39], spinnerTime, spinner.WithWriter(os.Stderr))
	s.Suffix = fmt.Sprintf(" %s", msg)
	s.FinalMSG = fmt.Sprintf(" %s\n", finalMsg)

	s.Start()

	return s
}
