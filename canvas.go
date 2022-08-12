package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Canvas struct {
	lastTimePainted time.Time
	buffer          string
}

func (c *Canvas) Close() {
	// flush last state
	c.Flush()
}

func NewCanvas() *Canvas {
	c := &Canvas{
		lastTimePainted: time.Now(),
	}
	return c
}

func (c *Canvas) Paint(s string) {
	c.buffer = s
	if time.Since(c.lastTimePainted) > 5*time.Second {
		clearTerminal()
		fmt.Println(s)
	}
}

func (c *Canvas) Flush() {
	clearTerminal()
	fmt.Println(c.buffer)
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
