package herder

import (
	"sync"
	"time"
)
type TaskConfig struct {
	Command string `json:"command"`
	Args []string `json:"args,omitempty"`
	MaxStdoutLen *int `json:"max_stdout_len,omitempty"`
	MaxStderrLen *int `json:"max_stderr_len,omitempty"`
}
type tasks struct {
	slice  []*task
	maxLen int
	m      sync.Mutex
}

func (ts *tasks) append(t *task) {
	if ts.maxLen > 0 {
		if len(ts.slice) == ts.maxLen {
			ts.m.Lock()
			newSlice := make([]*task, ts.maxLen, ts.maxLen)
			for i := 1; i < len(ts.slice); i++ {
				newSlice[i-1] = ts.slice[i]
			}
			newSlice[ts.maxLen-1] = t
			ts.slice = newSlice
			ts.m.Unlock()
		}
	} else {
		ts.slice = append(ts.slice, t)
	}
}

func (ts *tasks) remove(taskId int) {
	ts.m.Lock()
	defer ts.m.Unlock()
	for i := range ts.slice {
		if ts.slice[i] != nil && ts.slice[i].id == taskId {
			ts.slice = append(ts.slice[:i], ts.slice[i+1:]...)
			return
		}
	}
}

func (ts *tasks) getStates() []TaskState {
	ts.m.Lock()
	defer ts.m.Unlock()
	res := make([]TaskState, 0, len(ts.slice))
	for i := range ts.slice {
		res = append(res, ts.slice[i].getState())
	}
	return res
}

type task struct {
	id           int
	command      string
	args         []string
	startedAt    *time.Time
	finishedAt   *time.Time
	maxStdoutLen int
	maxStderrLen int
	p            *Process
}

type TaskState struct {
	TaskId     int        `json:"task_id"`
	Command    string     `json:"command"`
	Args       []string   `json:"args"`
	StdOut     []byte     `json:"std_out"`
	StdErr     []byte     `json:"std_err"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	ExitCode   *int       `json:"exit_code"`
	Error      error      `json:"error"`
}

func (t *task) getState() TaskState {
	var stdOut, stdErr []byte
	if t.p != nil && t.p.stdout != nil {
		stdOut = make([]byte, len(t.p.stdout.buffer), t.p.stdout.maxLen)
		copy(stdOut, t.p.stdout.buffer)
	}
	if t.p != nil && t.p.stderr != nil {
		stdErr = make([]byte, len(t.p.stderr.buffer), t.p.stdout.maxLen)
		copy(stdErr, t.p.stderr.buffer)
	}
	var startedAt, finishedAt *time.Time
	if t.startedAt != nil {
		startedAt = new(time.Time)
		*startedAt = *t.startedAt
	}
	if t.finishedAt != nil {
		finishedAt = new(time.Time)
		*finishedAt = *t.finishedAt
	}
	var args []string
	if len(t.args) != 0 {
		copy(args, t.args)
	}
	var exitCode *int
	if t.p != nil && t.p.exitCode != nil {
		exitCode = new(int)
		*exitCode = *t.p.exitCode
	}
	var err error
	if t.p != nil && t.p.err != nil {
		err = t.p.err
	}
	return TaskState{
		TaskId:     t.id,
		Command:    t.command,
		Args:       args,
		StdOut:     stdOut,
		StdErr:     stdErr,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		ExitCode:   exitCode,
		Error:      err,
	}
}

func (t *task) run() {
	t.p = newProcess(
		newOutputBuffer(t.maxStdoutLen),
		newOutputBuffer(t.maxStderrLen),
		t.command,
		t.args...,
	)
	t.startedAt = new(time.Time)
	*t.startedAt = time.Now()
	t.p.run()
	t.finishedAt = new(time.Time)
	*t.finishedAt = time.Now()
}

func (t *task) killProcess() error {
	if t.p != nil {
		return t.p.kill()
	}
	return nil
}
