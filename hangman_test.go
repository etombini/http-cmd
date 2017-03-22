package hangman_test

import (
	"github.com/etombini/hangman"
	"testing"
)

func TestLs(t *testing.T) {
	cmdline := "ls -la"
	h := hangman.Reaper(cmdline, 1)
	if h.Return_code != 0 {
		t.Error("Return code is not 0: ", h.Return_code)
	}
	if h.Stderr != "" {
		t.Error("There are errors on stderr")
	}
	if h.Stdout == "" {
		t.Error("There is no output on stdout")
	}
}

func TestOverTime(t *testing.T) {
	cmdline := "sleep 2"
	h := hangman.Reaper(cmdline, 1)
	if h.Return_code == 0 {
		t.Error("Return code is zero: ", h.Return_code)
	}
	if h.Stderr == "" {
		t.Error("There is no error on stderr")
	}
}
