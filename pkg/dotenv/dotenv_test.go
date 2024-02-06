package dotenv

import "testing"

func TestRead(t *testing.T) {
	envMap := map[string]string{
		"TEST1": "Test1",
		"TEST2": "Test2",
		"TEST3": "Test3",
	}
	envs, err := Read(mock)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}
	if len(envMap) != len(envs) {
		t.Errorf("length of maps not equal: got %v, want %v", len(envs), len(envMap))
	}
	for k, v := range envMap {
		if val, ok := envs[k]; !ok {
			t.Errorf("key %v expected but not found", k)
		} else {
			if v != val {
				t.Errorf("got %v, want %v", val, v)
			}
		}
	}
}
