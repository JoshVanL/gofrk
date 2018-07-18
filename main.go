package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
	}

	cmds, err := parse(os.Args[1:])
	if err != nil {
		printError(err)
	}

	if err := run(cmds); err != nil {
		printError(err)
	}

	os.Exit(0)
}

func run(cmds []*exec.Cmd) error {
	var wg sync.WaitGroup

	//fmt.Printf("%d\n", len(cmds))
	//for _, c := range cmds {
	//	fmt.Printf("%+v\n\n", c)
	//}

	for _, c := range cmds {
		wg.Add(1)
		//fmt.Printf("%+v\n\n", c)

		go func() {
			defer wg.Done()

			if err := c.Run(); err != nil {
				printCmdError(err)
				return
			}

			if err := c.Wait(); err != nil {
				printCmdError(err)
			}

			return
		}()
	}

	wg.Wait()

	return nil
}

func parse(args []string) ([]*exec.Cmd, error) {
	var cmds []*exec.Cmd
	var cmdArgs []string
	var result []error

	for _, a := range args {
		if !strings.Contains(a, ",") {
			cmdArgs = append(cmdArgs, a)
			continue
		}

		for _, s := range strings.Split(a, ",") {
			if s == "" {
				continue
			}
			cmdArgs = append(cmdArgs, s)

			cmd, err := createCmd(cmdArgs)
			if err != nil {
				result = append(result, err)
				continue
			}

			cmds = append(cmds, cmd)
			cmdArgs = []string{}
		}
	}

	if len(cmdArgs) > 0 {
		cmd, err := createCmd(cmdArgs)
		if err != nil {
			result = append(result, err)
		} else {
			cmds = append(cmds, cmd)
		}
	}

	return cmds, errorOrNil(result)
}

func createCmd(args []string) (*exec.Cmd, error) {
	if args == nil || len(args) == 0 {
		return nil, fmt.Errorf("command contains no arguments")
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &exec.Cmd{
		Args:   args,
		Dir:    dir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}, nil
}

func errorOrNil(result []error) error {
	if result == nil || len(result) == 0 {
		return nil
	}

	err := fmt.Errorf("a total of %d errors have occurred:\n", len(result))
	for _, e := range result {
		err = fmt.Errorf("%s\n* %s", err, e)
	}

	return err
}

func printHelp() {
	fmt.Printf(`Usage:
gofrk is a lightweight, no dependency forking utility written in Go.
Give comma separated terminal commands to be executed concurrently.
e.g. $ gofrk echo bar, sleep 300, touch hello world
`)
	os.Exit(0)
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "gofrk: fatal error: %s", err)
	os.Exit(1)
}

func printCmdError(err error) {
	fmt.Fprintf(os.Stderr, "gofrk: running a command went wrong: %s\n", err)
}
