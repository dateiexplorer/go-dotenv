package dotenv

import (
	"fmt"
	"strings"
	"testing"
)

var mock = &mockShell{}
var errParseLine = fmt.Errorf("string 'ERROR' is not allowed")

type mockShell struct{}

func (s *mockShell) Ignorable(line string) bool {
	return strings.HasPrefix(line, "#")
}

func (s *mockShell) ParseLine(line string) (string, string, error) {
	env := strings.SplitN(line, "=", 2)
	if strings.Contains(line, "ERROR") {
		return "", "", errParseLine
	}

	// Remove comments
	if idx := strings.IndexRune(env[1], '#'); idx > -1 {
		env[1] = strings.TrimSpace(env[1][:idx])
	}

	return env[0], env[1], nil
}

func TestIgnorable(t *testing.T) {
	input := map[string]bool{
		"TEST=Test":  false,
		"#TEST=Test": true,
		"Test=#Test": false,
	}
	for k, v := range input {
		ignorable := mock.Ignorable(k)

		if ignorable != v {
			t.Errorf("line %v: got %v, want %v", k, ignorable, v)
		}
	}
}

func TestParseLine(t *testing.T) {
	type output struct {
		k, v string
		err  error
	}
	input := map[string]*output{
		"TEST=Test":      {k: "TEST", v: "Test", err: nil},
		"TEST=#Test":     {k: "TEST", v: "", err: nil},
		"TEST=TestERROR": {k: "", v: "", err: errParseLine},
	}
	for line, out := range input {
		key, value, err := mock.ParseLine(line)
		if err != out.err {
			t.Errorf("line %v: got %v, want %v", line, err, out.err)
		}
		if key != out.k {
			t.Errorf("line %v: got %v, want %v", line, key, out.k)
		}
		if value != out.v {
			t.Errorf("line %v: got %v, want %v", line, value, out.v)
		}
	}
}
