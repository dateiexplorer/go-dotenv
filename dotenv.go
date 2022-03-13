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

	"github.com/dateiexplorer/go-dotenv/shells"
)

// Read reads all variables form a .env file into the environment. The function
// uses the Basic shell implementation, a cross platform Shell that implements
// a subset of the famout Bash shell.
// It returns a map which contains all successfully read environment variables
// from the file until an error occurred. The map keys are the variable names,
// the map values its correspondig values. Returns an error if any occurred.
func Read() (envs map[string]string, err error) {
	return ReadWith(shells.Basic, ".env")
}

// ReadWith reads all variables from a file which path is given in the env
// parameter into the environment. The function uses the syntax of the Shell
// shell.
// It returns a map which contains all successfully read environment variables
// from the file until an error occurred. The map keys are the variable names,
// the map values its correspondig values. Returns an error if any occurred.
//
// If you want to read from the default .env variable with a basic cross
// platform Shell, use the simple Read function.
func ReadWith(shell shells.Shell, env string) (envs map[string]string, err error) {
	envs = make(map[string]string)

	// Open the file form which the variables should be load.
	// By default this is the .env variable.
	f, err := os.Open(env)
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

// Load loads all variables defined in the .env file into the environment with
// a cross platform Basic Shell, which implements a subset of the functionality
// of the Bash shell.
//
// To get the variables use the functions from the os package of the standard
// library, e.g. os.Gentenv().
func Load() error {
	return LoadWith(shells.Basic, ".env")
}

// LoadWith loads all variables defined in the file env into the environment
// with the specific Shell shell.
//
// To get the variables use the functions from the os package of the standard
// library, e.g. os.Gentenv().
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
