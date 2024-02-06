package dotenv_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dateiexplorer/go-dotenv"
	"github.com/dateiexplorer/go-dotenv/shells"
)

var envMap = [][]string{
	{"TEST1", "Test1"},
	{"TEST2", "Test2"},
	{"TEST3", "Test3"},
}

func TestReadWith(t *testing.T) {
	envs, err := dotenv.ReadWith(shells.Basic, "testdata/.env")
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}
	if len(envMap) != len(envs) {
		t.Errorf("length of maps not equal: got %v, want %v", len(envs), len(envMap))
	}
	for _, env := range envMap {
		if val, ok := envs[env[0]]; !ok {
			t.Errorf("key %v expected but not found", env[0])
		} else {
			if env[1] != val {
				t.Errorf("got %v, want %v", val, env[1])
			}
		}
	}
}

func TestLoad(t *testing.T) {
	dotenv.LoadWith(shells.Basic, "testdata/.env")
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		v := strings.SplitN(env, "=", 2)
		envs[v[0]] = v[1]
	}
	for _, env := range envMap {
		if val, ok := envs[env[0]]; !ok {
			t.Errorf("key %v expected but not found", env[0])
		} else {
			if env[1] != val {
				t.Errorf("got %v, want %v", val, env[1])
			}
		}
	}
}
