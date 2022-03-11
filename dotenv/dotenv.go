package dotenv

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	if err := Load(); err != nil {
		log.Fatalln(err)
	}
}

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
	for i := 0; scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if shell.Ignorable(line) {
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

func Load() error {
	shell := new(BashShell)
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
