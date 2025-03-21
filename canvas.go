package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

type Canvas struct {
	clearWorks bool
	out        io.Writer

	mu     sync.Mutex
	buffer string
	errors []error
}

func NewCanvas(out io.Writer) *Canvas {
	c := &Canvas{
		out:    out,
		errors: []error{},
	}

	err := c.clear(out)
	if err == nil {
		c.clearWorks = true
	}

	return c
}

func (c *Canvas) Paint(table *SyncTable, clear bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buffer = table.Render()
	c.flush(clear)
	return nil
}

// Error collects error messages in order to print them at the end
func (c *Canvas) Error(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.errors = append(c.errors, err)
}

func (c *Canvas) Flush(clear bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flush(clear)
}

func (c *Canvas) flush(clear bool) {
	if clear {
		c.clearTerminal(c.out)
	}

	fmt.Fprint(c.out, c.buffer)

	for idx, err := range c.errors {
		fmt.Fprintln(c.out, idx, err)
	}
}

func (c *Canvas) Buffer() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buffer
}

func (c *Canvas) clearTerminal(out io.Writer) {
	if c.clearWorks {
		c.clear(out)
	}
}

func (c *Canvas) clear(out io.Writer) error {
	switch runtime.GOOS {
	case "darwin":
		return runCmd(out, "clear")
	case "linux":
		return runCmd(out, "clear")
	case "windows":
		return runCmd(out, "cmd", "/c", "cls")
	default:
		return runCmd(out, "clear")
	}
}

func runCmd(out io.Writer, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Env = os.Environ()
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}
