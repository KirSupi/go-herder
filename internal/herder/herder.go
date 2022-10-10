package herder

import "os/exec"

type Config struct {
}

type Herder struct {
	Processes []Process
	c         *Config
}

type Process struct {
	ID int

	Active bool
	Cmd    exec.Cmd
}

func New(c Config) *Herder {
	return &Herder{
		c: &c,
	}
}
