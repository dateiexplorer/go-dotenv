package shells

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// A basicShell is a cross platform shell which not depends on any operating
// system specific stuff.
// It provides a subset of the functionality and syntax of the famous bash
// shell and so is fully compatible.
type basicShell struct{}

func (b *basicShell) Ignorable(line string) bool {
	return strings.HasPrefix(line, "#")
}

func (b *basicShell) ParseLine(line string) (string, string, error) {
	// Split the environment variable by the first appereance of = in key
	// (name) and value part.
	// Parse them separately.
	env := strings.SplitN(line, "=", 2)
	key := parseKey(env[0])
	val := parseVal(env[1])

	return key, val, nil
}

var regExAllowedVarChars = regexp.MustCompile(`[_A-Za-z0-9]`)

type tokenScanner struct {
	input string
	head  int
}

func newTokenScanner(input string) *tokenScanner {
	return &tokenScanner{input: input, head: -1}
}

func (s *tokenScanner) Scan() bool {
	s.head++
	return s.head < len(s.input)
}

func (s *tokenScanner) StepBackwards() bool {
	if s.head > 0 {
		s.head--
		return true
	}
	return false
}

func (s *tokenScanner) Token() byte {
	return s.input[s.head]
}

func (s *tokenScanner) PreviousToken() (byte, error) {
	if s.head > 0 {
		return s.input[s.head-1], nil
	}

	return 0, fmt.Errorf("reached begin of input")
}

func parseKey(key string) string {
	return strings.TrimSpace(key)
}

func parseVal(val string) string {
	value := strings.Builder{}
	var scopeSingleQuoted bool
	var scopeDoubleQuoted bool

	for scanner := newTokenScanner(val); scanner.Scan(); {
		// Get the next Token from the scanner
		c := scanner.Token()

		// '\' escapes any character if its not single enquoted.
		// The '\' is not visible in the result, so skip it and write any
		// following character.
		if c == '\\' && !scopeSingleQuoted {
			scanner.Scan()
			value.WriteByte(scanner.Token())
			continue
		}

		// ''' starts and closes a a single enquoted string.
		if c == '\'' && !scopeDoubleQuoted {
			scopeSingleQuoted = !scopeSingleQuoted
			continue
		}

		// '"' starts and closes a double enquoted string.
		if c == '"' && !scopeSingleQuoted {
			scopeDoubleQuoted = !scopeDoubleQuoted
			continue
		}

		// '#' starts a comment if any whitespace char is detected before.
		if c == '#' {
			previous, err := scanner.PreviousToken()
			if err != nil {
				continue
			}

			r, _ := utf8.DecodeRune([]byte{previous})
			if unicode.IsSpace(r) {
				break
			}
		}

		// '$' starts a variable substitution if it not appears in a
		// single enquoted string.
		if c == '$' && !scopeSingleQuoted {
			variable := parseVar(scanner)
			value.WriteString(variable)
			continue
		}

		value.WriteByte(c)
	}

	return strings.TrimSpace(value.String())
}

func parseVar(scanner *tokenScanner) string {
	variable := strings.Builder{}
	for scanner.Scan() {
		c := scanner.Token()

		// var scopeBrackets bool
		if c == '{' {
			// scopeBrackets = true
			continue
		}

		if c == '}' {
			// scopeBrackets = false
			continue
		}

		if !regExAllowedVarChars.MatchString(string(c)) {
			break
		}

		variable.WriteByte(c)
	}

	// Step back to previous character for the next scan.
	scanner.StepBackwards()
	return os.Getenv(variable.String())
}
