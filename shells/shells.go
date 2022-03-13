// Package shells provides various implementations for shells to parse
// variables from a file into the environment.
// Further with the Shell interface this package gives the ability to implement
// individual interpreter on your own.
//
// A Shell in the context of the dotenv package is an execution environment
// that has the ability to parse variables as environment variables. A Shell
// in this context is nothing more than a set of specific syntax.
package shells

var (
	// Bash is an implementation of the famous Bash shell that is the default
	// shell in most Linux distributions.
	//
	// To parse variable in the environment this implementation uses the echo
	// command of a real bash shell. This assumes that a bash executable is
	// available under /bin/bash. So this implementation is only capable for
	// Linux environments.
	//
	// If you need a cross platform implementation use the Basic shell instead.
	// It provides a subset of the bash shell syntax and may reach out for the
	// most purposes.
	//
	// If you need fully support of all bash syntax you can use this Bash
	// implementation. The advantage using the bash shell instead of
	// implementing its whole syntax is that it ensures that all variable and
	// bash command substitutions working properly because the bash syntax can
	// be very complex.
	Bash = &bashShell{}

	// Basic is a cross platform shell which not depends on any operating
	// system specific stuff.
	// It provides a subset of the functionality and syntax of the famous
	// bash shell and so is fully compatible.
	//
	// Is most cases this is what you want and is should reach out for most
	// purposes.
	// But this implementation lacks some features such as bash command
	// substitutions.
	// If you would take care of those, probably you'll take a look at the Bash
	// implementation or create your own Shell that fit your needs.
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
