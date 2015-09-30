package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh/terminal"
)

func getPassword() string {
	fmt.Fprint(os.Stderr, "Password (typing will be hidden): ")
	p, err := terminal.ReadPassword(0)
	fmt.Println(string(p))
	fmt.Println(err)
	return "f"
	// This is real hacky, but it's the only way to get password reading working with
	// windows, unicode, git bash, cygwin, etc.
	script := `
	read -s password
	printf $password`
	exec.Command("stty", "-echo").Run()
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	fmt.Println(line)
	exec.Command("stty", "echo").Run()

	cmd := exec.Command("bash", "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	fmt.Fprintln(os.Stderr)
	ExitIfError(err)
	return string(output)
}
