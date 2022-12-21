package herder

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
)

type Process struct {
	ID           int
	Label        *string
	Command      string
	Params       string
	Active       bool
	Cmd          *exec.Cmd
	OutputBuffer *bytes.Buffer
}

type ProcessState struct {
	ID     int     `json:"id"`
	Label  *string `json:"label"`
	Active bool    `json:"active"`
	Output []byte  `json:"output"`
}

func (p *Process) getState() ProcessState {
	return ProcessState{
		ID:     p.ID,
		Label:  p.Label,
		Active: p.Active,
		Output: p.OutputBuffer.Bytes(),
	}
}
func (p *Process) run() error {
	cmdArgs := strings.Split(p.Command, " ")
	if len(cmdArgs) == 0 {
		return errors.New("empty command")
	}
	cmd := cmdArgs[0]
	args := make([]string, 0, len(cmdArgs)-1)
	for _, s := range cmdArgs[1:] {
		args = append(args, s)
	}
	for _, s := range strings.Split(p.Params, " ") {
		args = append(args, s)
	}
	log.Println("CMD:", cmd)
	p.Cmd = exec.Command(cmd, args...)

	log.Println("run", p.ID, p.Params, p.Cmd.Path)
	if p.Active || p.Cmd != nil || p.Cmd.Process != nil {
		return errors.New("process already running")
	}
	p.Active = true
	p.OutputBuffer = bytes.NewBuffer(nil)

	p.Cmd.Stdout = p.OutputBuffer
	p.Cmd.Stderr = p.OutputBuffer
	go func() {
		if err := p.Cmd.Run(); err != nil {
			log.Printf("process #%d ends with error: %s\n", p.ID, err.Error())
		} else {
			log.Printf("process #%d ends\n", p.ID)
		}
		p.Active = false
	}()
	return nil
}
func (p *Process) kill() error {
	p.Active = false
	if p.Cmd != nil {
		if p.Cmd.Process != nil {
			err := p.Cmd.Process.Kill()
			p.Cmd = nil
			return err
		}
	}
	return nil
}
