package main

import (
	"bytes"
	"testing"
)

func TestCLI(t *testing.T) {
	for _, tc := range []struct {
		args []string
		code int
		out  string
	}{
		{[]string{"suggest", "pwd"}, 0, "pwd\n"},
		{[]string{"suggest", "delete x"}, 1, ""},
		{[]string{}, 2, ""},
		{[]string{"suggest", "pwd", "whoami"}, 2, ""},
	} {
		var out, errs bytes.Buffer
		code := run(tc.args, &out, &errs)
		if code != tc.code || out.String() != tc.out {
			t.Fatalf("%v: %d %q %q", tc.args, code, out.String(), errs.String())
		}
		if code != 0 && errs.Len() == 0 {
			t.Fatal("missing error")
		}
	}
}

func TestReleaseCLI(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"catalog"}, {"config"}, {"version"}, {"init", "bash"}} {
		var out, errs bytes.Buffer
		if run(args, &out, &errs) != 0 || out.Len() == 0 || errs.Len() != 0 {
			t.Fatalf("failed %v: %s", args, errs.String())
		}
	}
	for _, args := range [][]string{{"init", "zsh"}, {"config", "auto_execute=true"}, {"doctor", "extra"}, {"catalog", "extra"}} {
		var out, errs bytes.Buffer
		if run(args, &out, &errs) == 0 || out.Len() != 0 {
			t.Fatalf("accepted %v", args)
		}
	}
}
