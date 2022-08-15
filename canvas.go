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
	errors          []error
}

func (c *Canvas) Close() {
	// flush last state
	c.Flush()
}

func NewCanvas() *Canvas {
	c := &Canvas{
		lastTimePainted: time.Now(),
		errors:          []error{},
	}
	return c
}

func (c *Canvas) Paint(s string) {
	c.buffer = s
	if time.Since(c.lastTimePainted) > 5*time.Second {
		clearTerminal()
		fmt.Println(s)
		c.lastTimePainted = time.Now()
	}
}

// Error collects error messages in order to print them at the end
func (c *Canvas) Error(err error) {
	c.errors = append(c.errors, err)
}

func (c *Canvas) Flush() {
	clearTerminal()
	fmt.Println(c.buffer)

	for idx, err := range c.errors {
		fmt.Fprintln(os.Stderr, idx, err)
	}
}

func (c *Canvas) Save(filePath string) error {
	fmt.Println(c.buffer)

	return os.WriteFile(filePath, []byte(c.buffer), 0700)
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
