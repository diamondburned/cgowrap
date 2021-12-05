package shortflag

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	in := []string{
		"argument1",
		"-v", "1",
		"-v2",
		"argument2",
		"--value", "3",
		"--short",
		"--unknown", "-u",
		"argument3",
	}

	f, err := Parse(in, Opts{
		ValueFlags: []string{"-v", "-v2", "--value"},
		BlankFlags: []string{"--short"},
	})
	if err != nil {
		t.Fatal(err)
	}

	expect := &Flags{
		Args: []string{"argument1", "argument2", "--unknown", "-u", "argument3"},
		Flags: map[string][]string{
			"--short": nil,
			"--value": {"3"},
			"-v":      {"1", "2"},
		},
	}

	if !reflect.DeepEqual(expect, f) {
		t.Errorf("expected: %#v", expect)
		t.Errorf("got:      %#v", f)
	}
}
