package install

import (
	"fmt"
	"io"
)

func SetupOSXDefaults(output io.Writer) {
	if !runOSXDefaults(output) {
		return
	}
	fmt.Fprintln(output, "OSX defaults setup complete.")
}

func runOSXDefaults(output io.Writer) bool {
	if err := applyOSXDefaults(); err != nil {
		fmt.Fprintf(output, "Failed to apply OSX defaults: %v\n", err)
		return false
	}
	return true
}

func applyOSXDefaults() error {
	return runCommand("bash", "-c", "source ./osxdefaults.sh")
}
