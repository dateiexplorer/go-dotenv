package dotenv

import (
	"fmt"
	"os/exec"
	"strings"
)

type Shell interface {
	Ignorable(line string) bool
	ParseLine(line string) (key string, value string, err error)
}

type BashShell struct{}

func (s *BashShell) Ignorable(line string) bool {
	return strings.HasPrefix(line, "#")
}

func (s *BashShell) ParseLine(line string) (string, string, error) {
	// Execute the comment
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("echo -n %v", line))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to echo: %v desc: %v", err, string(out))
	}

	// Split result in key and value
	env := strings.SplitN(string(out), "=", 2)
	return env[0], env[1], nil
}
