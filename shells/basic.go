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
// system specific stuff or installed executables.
// It provides a subset of the functionality and syntax of the famous Bash
// shell and so is fully compatible.
//
// To parse variables this implementation uses a tokenScanner which read the
// variable char by char and has the ability to jump forward and backward in
// the string.
type basicShell struct{}

// Ignorable returns true if a line starting with a # sign. This is a comment
// in the Bash syntax.
func (b *basicShell) Ignorable(line string) bool {
	return strings.HasPrefix(line, "#")
}

// ParseLine consumes a whole line and parses it into key and value.
// Returns an error if parsing failed, e.g. because of syntax errors.
func (b *basicShell) ParseLine(line string) (string, string, error) {
	// Split the environment variable by the first appereance of = in key
	// (name) and value part.
	// Parse them separately.
	env := strings.SplitN(line, "=", 2)
	key := parseKey(env[0])
	val := parseVal(env[1])

	return key, val, nil
}

// regExAllowedVarChars is a regular expression that holds all allowed
// characters in a variable name.
var regExAllowedVarChars = regexp.MustCompile(`[_A-Za-z0-9]`)

// tokenScanner is a scanner that has the ability to jump forward and
// backwards in the input string.
type tokenScanner struct {
	input string // input sequence to scan
	head  int    // current scan position
}

// newTokenScanner returns a pointer to a new tokenScanner.
// Need the input sequence.
//
// To get first token call the Scan function.
func newTokenScanner(input string) *tokenScanner {
	return &tokenScanner{input: input, head: -1}
}

// scan set the head to the next char of the input.
// If the end of the input is reached this function returns false.
func (s *tokenScanner) scan() bool {
	s.head++
	return s.head < len(s.input)
}

// stepBackwards set the head to the previous char of the input.
// If the begin of the input is reached this function returns false.
func (s *tokenScanner) stepBackwards() bool {
	s.head--
	return s.head > 0
}

// token returns the char at the current head position.
func (s *tokenScanner) token() byte {
	return s.input[s.head]
}

// previousToken returns the char at the previous position.
// If there is no previous position the function returns an error
// and a '0' byte.
func (s *tokenScanner) previousToken() (byte, error) {
	if s.head > 0 {
		return s.input[s.head-1], nil
	}

	return 0, fmt.Errorf("reached begin of input")
}

// parseKey parses a environment vairable name and returns it.
func parseKey(key string) string {
	key = strings.TrimPrefix(key, "export")
	return strings.TrimSpace(key)
}

// parseVal parses a environment variables value and returns it.
func parseVal(val string) string {
	value := strings.Builder{}
	var scopeSingleQuoted bool
	var scopeDoubleQuoted bool

	for scanner := newTokenScanner(val); scanner.scan(); {
		// Get the next Token from the scanner
		c := scanner.token()

		// '\' escapes any character if its not single enquoted.
		// The '\' is not visible in the result, so skip it and write any
		// following character.
		if c == '\\' && !scopeSingleQuoted {
			scanner.scan()
			value.WriteByte(scanner.token())
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
			previous, err := scanner.previousToken()
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

// parseVar extends the environment variable in an environment variable value.
// Returns the extended environment variable.
func parseVar(scanner *tokenScanner) string {
	variable := strings.Builder{}
	for scanner.scan() {
		c := scanner.token()

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
	scanner.stepBackwards()
	return os.Getenv(variable.String())
}
