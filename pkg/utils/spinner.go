package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

/* StartSpinner starts a spinner with a message and then replaces it with a done message
 * when the spinner is stopped.
 */
func StartSpinner(msg string, doneMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[39], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = fmt.Sprintf(" %s", msg)
	s.FinalMSG = fmt.Sprintf(" %s\n", doneMsg)

	s.Start()

	return s
}
