package herder

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

type Process struct {
	ID         int
	Label      *string
	Command    string
	Params     string
	Active     bool
	Cmd        *exec.Cmd
	outputPipe io.ReadCloser
}

type ProcessState struct {
	ID     int     `json:"id"`
	Label  *string `json:"label"`
	Active bool    `json:"active"`
	Output []byte  `json:"output"`
}

func (p *Process) getState() ProcessState {
	scanner := bufio.NewScanner(p.outputPipe)
	if scanner.Err() != nil {
		log.Println(p.ID, "scanner.Err():", scanner.Err())
	}
	for scanner.Scan() {
		log.Println(p.ID, "buf:", scanner.Bytes(), string(scanner.Bytes()))
	}
	return ProcessState{
		ID:     p.ID,
		Label:  p.Label,
		Active: p.Active,
		Output: scanner.Bytes(),
	}
}
func (p *Process) run() error {
	log.Println(p.ID, *p.Label, p.Params, p.Cmd.Path)
	if p.Cmd == nil {
		return nil
	}
	if p.Cmd.Process != nil || p.Active {
		return errors.New("process already running")
	}
	p.Active = true
	pr, pw := io.Pipe()
	p.Cmd.Stdout = pw
	go func() {
		go func() {
			if _, err := io.Copy(os.Stdout, pr); err != nil {
				log.Fatal(err)
			}
			err := pr.Close()
			if err != nil {
				log.Println("PipeWriter Error:", err.Error())
			}
		}()
		if err := p.Cmd.Run(); err != nil {
			log.Printf("process #%d ends with error: %s\n", p.ID, err.Error())
		} else {
			log.Printf("process #%d ends\n", p.ID)
		}
		//if output, err := p.Cmd.CombinedOutput(); err != nil {
		//	log.Printf("process #%d output: %s\n", p.ID, string(output))
		//	log.Printf("process #%d ends with error: %s\n", p.ID, err.Error())
		//} else {
		//	log.Printf("process #%d output: %s\n", p.ID, string(output))
		//	log.Printf("process #%d ends\n", p.ID)
		//}
		p.Active = false

		err := pw.Close()
		if err != nil {
			log.Println("PipeWriter Error:", err.Error())
		}
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
