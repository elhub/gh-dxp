package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

/* StartSpinner starts a spinner with a message and a done message
 */
func StartSpinner(msg string, doneMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s", msg)
	s.FinalMSG = fmt.Sprintf(" %s\n", doneMsg)

	s.start()

	return s
}
