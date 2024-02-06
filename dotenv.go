// Package dotenv provides functions to load variables form a file into the
// environment. This package is heavily influcenced by the original Ruby
// implementation (https://github.com/bkeepers/dotenv) and various other
// Go libraries, e.g. https://github.com/joho/godotenv.
//
// Other than the most existing libraries this module provides a well
// defined interface giving the ability to implement further interpreters (so
// called 'Shells') to support individual syntax.
//
// Provides functionality to load variables defined in a file (normaly
// '.env'), e.g.
//
//  SOME_ENV_VAR=somevalue
//
// into the environment by calling
//
//  dotenv.Load()
//
// Afterwards the variables will be available through standard library
// functions, such as
//
//  os.Gentenv("SOME_ENV_VAR").
//
package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dateiexplorer/go-dotenv/shells"
)

// Read is a special version of ReadWith and uses the BasicShell implementation
// to load variables from a file ".env" into a map.
// It returns a map which contains all successfully read environment variables
// from the file until an error occured. The map keys are the variable names,
// the map values their corresponding values. Returns an error if any occurred.
//
// The BasicShell is a cross platform Shell that implements a subset of the
// famous Bash syntax. This should be sufficient for the most use cases.
func Read() (envs map[string]string, err error) {
	return ReadWith(shells.Basic, ".env")
}

// ReadWith reads all variables from a file which path is given in the env
// parameter into a map in form map[env_name]env_value. The function uses the
// syntax of the Shell shell to interpret the file.
// It returns a map which contains all successfully read environment variables
// from the file until an error occurred. The map keys are the variable names,
// the map values their corresponding values. Returns an error if any occurred.
//
// If you want to read from the standard ".env" file with a basic cross
// platform Shell that imitates the Bash syntax, use Read instead.
//
// If you want to load the variables directly into the environment use the
// LoadWith function.
func ReadWith(shell shells.Shell, env string) (envs map[string]string, err error) {
	envs = make(map[string]string)

	// Open the file from which the variables should be load.
	// By default this is the ".env" file.
	f, err := os.Open(env)
	if err != nil {
		return envs, fmt.Errorf("cannot open env file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Trim whitespace
		line := strings.TrimSpace(scanner.Text())
		// Skip empty and other ignorable files
		if len(line) == 0 || shell.Ignorable(line) {
			continue
		}

		key, value, err := shell.ParseLine(line)
		if err != nil {
			return envs, fmt.Errorf("cannot parse line %v: %w", line, err)
		}

		envs[key] = value
	}

	return envs, nil
}

// Load is a special version of LoadWith and uses the BasicShell implementation
// to load variables from a file ".env" into the environment.
// Returns an error if any occurred.
//
// The BasicShell is a cross platform Shell that implements a subset of the
// famous Bash syntax. This should be sufficient for the most use cases.
func Load() error {
	return LoadWith(shells.Basic, ".env")
}

// LoadWith loads all variables defined in the env file into the environment
// with the specific syntax of a Shell shell.
// Returns an error if any occurred.
//
// To get the variables use the functions from the os package of the standard
// library, e.g. os.Gentenv("SOME_VAR_ENV").
//
// If you want to load from the standard ".env" file with a basic cross
// platform Shell that imitates the Bash syntax, use Load instead.
func LoadWith(shell shells.Shell, env string) error {
	envs, err := ReadWith(shell, env)
	if err != nil {
		return err
	}

	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf(
				"cannot add var with key '%v' and value '%v' to the environment: %w", k, v, err)
		}
	}

	return nil
}
