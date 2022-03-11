package dotenv

import (
	"fmt"
	"os/exec"
	"strings"
)

// Shell is the interface that wraps an execution environment for parsing env
// variables for various shells and operating systems.
type Shell interface {
	// Ignorable returns true if the line should be ignored and continue with
	// the next line.
	// This is useful to ignore comments in various shell script languages.
	// Empty lines are always skipped and may not take into account here.
	Ignorable(line string) bool

	// ParseLine parses a single line of the .env file into an environment
	// variable. The line is already trimmed and doesn't start or end with any
	// whitespace character.
	// This function should return the key (or name) of the env var and its
	// value.
	// If an error occured during parsing the error should be returned in the
	// err var and the key and value should be an empty string.
	ParseLine(line string) (key string, value string, err error)
}

// A BashShell is an implementation of Shell for the famous bash shell which is
// in most Linux environments the standard shell.
type BashShell struct{}

// Ignorable support bash comments in an environment file which started with a
// '#' symbol.
func (s *BashShell) Ignorable(line string) bool {
	return strings.HasPrefix(line, "#")
}

// ParseLine parses a single environment variable using an echo command in the
// bash shell. This assumes that bash executable is available under /bin/bash.
//
// The advantage of using the bash shell instead of implementing its whole
// syntax is that all variable and bash command substitutions working properly.
// Bash syntax is very complex. If a other implementation is needed you can
// create it on your own by creating a new type implements the Shell interface.
//
// It returns the name and value of the the environment variable or empty
// strings if an error occured during the parsing. Then the error variable is
// set. Otherwise the error is nil.
//
// This implementation supports the 'export' keyword for variables, so that a
// .env file can also be sourced in a bash environment if needed.
func (s *BashShell) ParseLine(line string) (string, string, error) {
	// Echo the line in the bash shell. Use the output to create key value
	// pair.
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("echo -n %v", line))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to echo: %v desc: %v", err, string(out))
	}

	// Split result in key and value
	env := strings.SplitN(string(out), "=", 2)

	// Remove the 'export' keyword in name variable, trim whitespace if
	// necessary.
	env[0] = strings.TrimSpace(strings.TrimPrefix(env[0], "export"))
	return env[0], env[1], nil
}
