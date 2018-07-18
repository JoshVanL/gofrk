package main

import (
	"os/exec"
	"strings"
	"testing"
)

func Test_Parse(t *testing.T) {

	testArgs := strings.Split(
		"foo bar, foo bar",
		" ",
	)

	cmds, err := parse(testArgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
		return
	}

	mustLength(2, len(cmds), t)
	mustArgs([]string{"foo", "bar"}, cmds[0], t)
	mustArgs([]string{"foo", "bar"}, cmds[1], t)

	testArgs = strings.Split(
		",foo,bar,koo,boo,",
		" ",
	)
	cmds, err = parse(testArgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
		return
	}

	mustLength(4, len(cmds), t)
	mustArgs([]string{"foo"}, cmds[0], t)
	mustArgs([]string{"bar"}, cmds[1], t)
	mustArgs([]string{"koo"}, cmds[2], t)
	mustArgs([]string{"boo"}, cmds[3], t)

	testArgs = strings.Split(
		",,foo,,, bar \"foo bar\"",
		" ",
	)
	cmds, err = parse(testArgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
		return
	}

	mustLength(2, len(cmds), t)
	mustArgs([]string{"foo"}, cmds[0], t)
	mustArgs([]string{"bar", "\"foo", "bar\""}, cmds[1], t)
}

func mustLength(exp, got int, t *testing.T) {
	if exp != got {
		t.Fatalf("wrong number of commands, exp=%d got=%d", exp, got)
	}
}

func mustArgs(args []string, cmd *exec.Cmd, t *testing.T) {
	if !matchSlice(args, cmd.Args) {
		t.Errorf("command has incorrect arguments, exp=%s, got=%s", args, cmd.Args)
	}
}

func matchSlice(a, b []string) bool {
	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i, c := range a {
		if b[i] != c {
			return false
		}
	}

	return true
}
