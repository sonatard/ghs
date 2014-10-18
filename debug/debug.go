package debug

import (
	"fmt"
)

func Printf(format string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(format, args...)
	}
}
