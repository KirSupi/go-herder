package herder

import (
	"errors"
	"log"
	"os/exec"
)

type Process struct {
	ID     int
	Label  *string
	Params string
	Active bool
	Cmd    *exec.Cmd
}

type ProcessState struct {
	ID     int     `json:"id"`
	Label  *string `json:"label"`
	Active bool    `json:"active"`
}

func (p *Process) getState() ProcessState {
	return ProcessState{
		ID:     p.ID,
		Label:  p.Label,
		Active: p.Active,
	}
}
func (p *Process) run() error {
	if p.Cmd == nil {
		return nil
	}
	if p.Cmd.Process != nil || p.Active {
		return errors.New("process already running")
	}
	p.Active = true
	go func() {
		if output, err := p.Cmd.CombinedOutput(); err != nil {
			log.Printf("process #%d output: %s\n", p.ID, string(output))
			log.Printf("process #%d ends with error: %s\n", p.ID, err.Error())
		} else {
			log.Printf("process #%d output: %s\n", p.ID, string(output))
			log.Printf("process #%d ends\n", p.ID)
		}
		p.Active = false
	}()
	return nil
}

func (p *Process) kill() error {
	if p.Cmd != nil {
		if p.Cmd.Process != nil {
			return p.Cmd.Process.Kill()
		}
	}
	return nil
}
