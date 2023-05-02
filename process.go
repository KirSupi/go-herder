package herder

import (
	"io"
	"os/exec"
)

type Process struct {
	cmd      *exec.Cmd
	stdout   *OutputBuffer
	stderr   *OutputBuffer
	stdin    *io.Writer
	exitCode *int
	err      error
}

func newProcess(stdout *OutputBuffer, stderr *OutputBuffer, command string, args ...string) *Process {
	p := &Process{
		cmd:    exec.Command(command, args...),
		stdout: stdout,
		stderr: stderr,
	}
	return p
}

func (p *Process) run() {
	if p.cmd == nil {
		return
	}
	if p.stdout != nil {
		p.cmd.Stdout = p.stdout
	}
	if p.stderr != nil {
		p.cmd.Stderr = p.stderr
	}
	p.err = p.cmd.Run()
	if p.cmd.ProcessState != nil {
		p.exitCode = new(int)
		*p.exitCode = p.cmd.ProcessState.ExitCode()
	}
}

func (p *Process) kill() error {
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}
