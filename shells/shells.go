// Package shells provides various implementations of shell syntax to parse
// variables from a file into the environment.
// Further with the Shell interface this package gives the ability to implement
// individual interpreter on your own.
//
// A Shell in the context of the dotenv package is an execution environment
// that has the ability to parse variables as environment variables. A Shell
// in this context is nothing more than a set of specific syntax to interpret
// environment variables.
//
package shells

var (
	// Bash is an implementation of the famous Bash shell that is the default
	// shell in most Linux distributions.
	//
	// To parse environment variables this implementation uses the "echo"
	// command of the real Bash. This assumes that a Bash executable is
	// available under "/bin/bash". So this implementation is only capable for
	// Linux environments.
	//
	// If you need a cross platform implementation use the Basic shell instead.
	// It provides a subset of the Bash syntax and may be sufficient for the
	// most purposes.
	//
	// If you need fully support of all Bash syntax you can use this Bash
	// implementation. The advantage using a real Bash instead of implementing
	// its whole syntax is that it ensures that all variable and Bash command
	// substitutions working properly because the bash syntax can be very
	// complex.
	//
	// This implementation supports command substitutions, e.g.
	//
	//  SOME_ENV_VAR=/some/path/$(whoami)
	//
	Bash = &bashShell{}

	// Basic is a cross platform shell which not depends on any operating
	// system specific stuff or the availablity of executables.
	// It provides a subset of the functionality and syntax of the famous Bash
	// and so it is fully compatible.
	//
	// Is most cases this is what you want and it should be sufficient for most
	// purposes.
	// But this implementation lacks some features such as Bash command
	// substitutions.
	// If you would take care of those, probably you'll take a look at the Bash
	// implementation or create your own Shell implementation that fit your
	// needs.
	Basic = &basicShell{}
)

// Shell is the interface that wraps an execution environment for parsing
// environment variables for various shell syntax and operating systems.
type Shell interface {
	// Ignorable returns true if the line should be ignored and continue with
	// the next line.
	// This is useful to ignore comments in various shell script languages.
	// Empty lines are always skipped and may not take into account here.
	Ignorable(line string) bool

	// ParseLine parses a single line of the .env file into an environment
	// variable. The line is already trimmed and doesn't start or end with any
	// whitespace character.
	// This function should return the key (name) of the variable and its
	// value.
	// If an error occured during parsing the error should be returned in the
	// err variable and the key and value should be an empty string by
	// convention.
	ParseLine(line string) (key string, value string, err error)
}
