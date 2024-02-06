// Package dotenv provides functions to use variables form an .env file int the
// current environment.
//
// Use the standard library functions to use the loaded env vars, e.g.
// os.Gentenv().
package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Read reads all environment variables form the .env file with the syntax
// parsing of a specific Shell shell.
// The function returns a map of parsed environment variables. Names are
// map keys and the values are the correspoding map values. If any error
// occured the function returns it. The map contains all env variables
// that were successfully read until the error occurres.
func Read(shell Shell) (envs map[string]string, err error) {
	envs = make(map[string]string)

	// Open the file form which the variables should be load.
	// By default this is the .env variable.
	f, err := os.Open(".env")
	if err != nil {
		return envs, fmt.Errorf("cannot open .env file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Trim whitespace.
		line := strings.TrimSpace(scanner.Text())
		// Skip empty and other ignorable files.
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

// Load loads all variables defined in the environment file into the current
// environment with the os specific default shell. For Linux this is the bash
// shell. Returns an error if any occurres.
//
// To get the variables use the functions from the os package of the standard
// library, e.g. os.Gentenv().
func Load() error {
	shell := new(BashShell)

	// Read all environment variables from the .env file into a map.
	envs, err := Read(shell)
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
