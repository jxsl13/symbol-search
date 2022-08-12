package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/gdamore/tcell/v2/encoding"
)

func init() {
	encoding.Register()
}

type Canvas struct {
}

func (c *Canvas) Close() {
}

func NewCanvas() *Canvas {
	c := &Canvas{}
	return c
}

func (c *Canvas) Paint(s string) {
	clearTerminal()
	fmt.Println(s)
}

func clearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
