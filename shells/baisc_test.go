package shells

import (
	"os"
	"testing"
)

func TestBasicIgnore(t *testing.T) {
	input := map[string]bool{
		"#TEST=Test":          true,
		"#!/bin/bash":         true,
		"TEST=Test":           false,
		"TEST=Test # Comment": false,
	}
	for k, v := range input {
		ignorable := Basic.Ignorable(k)
		if ignorable != v {
			t.Errorf("line %v: got %v, want %v", k, ignorable, v)
		}
	}
}

func TestBasicParseLineBasic(t *testing.T) {
	input := map[string]*output{
		// Unquoted most basic variable
		"TEST=Test": {k: "TEST", v: "Test", err: nil},
		// Ignore whitespace
		"  \t TEST=  Test \t": {k: "TEST", v: "Test", err: nil},
		// Allow empty variables
		"TEST=": {k: "TEST", v: "", err: nil},
	}
	testBasicParseLine(t, input)
}

func TestBasicParseLineComments(t *testing.T) {
	input := map[string]*output{
		// Ignore comments with the following expections
		// -> # is in single quotes
		// -> # is escaped with \
		"TEST=Test # Ignore comment":       {k: "TEST", v: "Test", err: nil},
		"TEST='#NoComment'":                {k: "TEST", v: "#NoComment", err: nil},
		"TEST = \"#NoComment\"":            {k: "TEST", v: "#NoComment", err: nil},
		"TEST=\\#NoComment":                {k: "TEST", v: "#NoComment", err: nil},
		"TEST='Test' # Comment":            {k: "TEST", v: "Test", err: nil},
		"TEST=\"Test1\"\"Test2\" #Comment": {k: "TEST", v: "Test1Test2", err: nil},
	}
	testBasicParseLine(t, input)
}

func TestBasicParseLineVarSubstitution(t *testing.T) {
	os.Clearenv()
	os.Setenv("HOME", "/home/me")
	input := map[string]*output{
		// Substitute variables starting with $ with the following rules
		// -> Do not substitute if $ is in single quotes
		// -> Do not substitute if $ is escaped with \
		"TEST=$HOME":      {k: "TEST", v: "/home/me", err: nil},
		"TEST=\"$HOME\"":  {k: "TEST", v: "/home/me", err: nil},
		"TEST='$HOME'":    {k: "TEST", v: "$HOME", err: nil},
		"TEST=\\$HOME":    {k: "TEST", v: "$HOME", err: nil},
		"TEST=${HOME}":    {k: "TEST", v: "/home/me", err: nil},
		"TEST=$HOME/test": {k: "TEST", v: "/home/me/test", err: nil},
	}
	testBasicParseLine(t, input)
}

type output struct {
	k, v string
	err  error
}

func testBasicParseLine(t *testing.T, input map[string]*output) {
	for line, out := range input {
		key, value, err := Basic.ParseLine(line)
		if err != out.err {
			t.Errorf("line %v (err): got %v, want %v", line, err, out.err)
		}
		if key != out.k {
			t.Errorf("line %v (key): got %v, want %v", line, key, out.k)
		}
		if value != out.v {
			t.Errorf("line %v (val): got %v, want %v", line, value, out.v)
		}
	}
}
