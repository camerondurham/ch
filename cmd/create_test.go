package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	cmd := rootCmd

	args := []string{"create", "test1"}
	cmd.SetArgs(args)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	expectsContains := "requires at least 1 arg"
	if strings.Contains(string(out), expectsContains) == false {
		t.Fatalf("expected to contain: %s, got: %s", expectsContains, string(out))
	} else {
		t.Log("correct error")
	}

}
