package logging

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

// This output is shown if a panic happens.
const panicOutput = `
!!!!!!!!!!!!!!!!!!!!!!!!!!! OPSOS CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!
opsos crashed! This is always indicative of a bug within opsos.
Please report the crash with Opsos[1] so that we can fix this.
When reporting bugs, please include your opsos version, the stack trace
shown below, and any additional information which may help replicate the issue.
[1]: https://github.com/neermitt/opsos/issues
!!!!!!!!!!!!!!!!!!!!!!!!!!! OPSOS CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!
`

// In case multiple goroutines panic concurrently, ensure only the first one
// recovered by PanicHandler starts printing.
var panicMutex sync.Mutex

// PanicHandler is called to recover from an internal panic in Terraform, and
// augments the standard stack trace with a more user-friendly error message.
// PanicHandler must be called as a deferred function, and must be the first
// defer called at the start of a new goroutine.
func PanicHandler() {
	// Have all managed goroutines checkin here, and prevent them from exiting
	// if there's a panic in progress. While this can't lock the entire runtime
	// to block progress, we can prevent some cases where Terraform may return
	// early before the panic has been printed out.
	panicMutex.Lock()
	defer panicMutex.Unlock()

	recovered := recover()
	if recovered == nil {
		return
	}

	fmt.Fprint(os.Stderr, panicOutput)
	fmt.Fprint(os.Stderr, recovered, "\n")

	// When called from a deferred function, debug.PrintStack will include the
	// full stack from the point of the pending panic.
	debug.PrintStack()

	// An exit code of 11 keeps us out of the way of the detailed exitcodes
	// from plan, and also happens to be the same code as SIGSEGV which is
	// roughly the same type of condition that causes most panics.
	os.Exit(11)
}
