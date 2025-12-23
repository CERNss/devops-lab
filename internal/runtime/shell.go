package runtime

import (
	"io"
	"os"
	"os/exec"
)

type ShellRuntime struct{}

func (sr *ShellRuntime) Run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (sr *ShellRuntime) RunQuiet(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = io.Discard
	c.Stderr = io.Discard
	return c.Run()
}

func RunShell(script string, env map[string]string) error {
	cmd := exec.Command("bash", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	return cmd.Run()
}

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
