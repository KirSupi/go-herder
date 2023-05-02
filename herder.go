package herder

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Herder struct {
	queue           TasksQueue
	finished        tasks
	active          tasks
	lastTaskId      int
	maxWorkersCount int
	workers         chan int
	m               sync.Mutex
	logger          *Logger
	defaultMaxStdoutLen int
	defaultMaxStderrLen int
}

type Config struct {
	MaxWorkersCount int `json:"max_workers_count,omitempty"`
	Logger          *Logger
	DefaultMaxStdoutLen int `json:"default_max_stdout_len"`
	DefaultMaxStderrLen int `json:"default_max_stderr_len"`
}

func New(c Config) *Herder {
	return &Herder{
		queue: TasksQueue{
			ch: make(chan *task),
		},
		maxWorkersCount: c.MaxWorkersCount,
		logger:          c.Logger,
		defaultMaxStdoutLen: c.DefaultMaxStdoutLen,
		defaultMaxStderrLen: c.DefaultMaxStderrLen,
	}
}

func (h *Herder) AddToQueue(tc TaskConfig) (taskId int) {
	h.lastTaskId++
	if tc.MaxStdoutLen == nil {
		tc.MaxStdoutLen = new(int)
		*tc.MaxStdoutLen = h.defaultMaxStdoutLen
	}
	if tc.MaxStderrLen == nil {
		tc.MaxStderrLen = new(int)
		*tc.MaxStderrLen = h.defaultMaxStderrLen
	}
	t := &task{
		id:      h.lastTaskId,
		command: tc.Command,
		args:    tc.Args,
		maxStdoutLen: *tc.MaxStdoutLen,
		maxStderrLen: *tc.MaxStderrLen,
	}
	h.queue.add(t)
	return t.id
}
func (h *Herder) RemoveFromQueue(taskId int) error {
	return h.queue.remove(taskId)
}
func (h *Herder) ClearQueue() {
	h.queue.clear()
}

func (h *Herder) Run() {
	if h.maxWorkersCount > 0 {
		h.workers = make(chan int, h.maxWorkersCount)
		for i := 1; i <= h.maxWorkersCount; i++ {
			go h.worker(i, nil)
		}
		for i := 1; i <= h.maxWorkersCount; i++ {
			<-h.workers
		}
		close(h.workers)
		h.workers = nil
	} else {
		var i int
		for t := range h.queue.ch {
			i++
			go h.worker(i, t)
		}
	}
}

func (h *Herder) log(v ...any) {
	if h.logger != nil {
		ts := time.Now().Format(time.DateTime)
		v = append([]any{ts}, v...)
		(*h.logger).Println(v...)
	}
}

func (h *Herder) worker(workerId int, t *task) {
	h.log(fmt.Sprintf("worker #%d started", workerId))
	if t == nil {
		for t = range h.queue.ch {
			h.runTask(workerId, t)
		}
	} else {
		h.runTask(workerId, t)
	}
	h.log(fmt.Sprintf("worker #%d stopped", workerId))
}

func (h *Herder) runTask(workerId int, t *task) {
	h.log(fmt.Sprintf("worker #%d run task #%d", workerId, t.id))
	h.active.append(t)
	t.run()
	h.active.remove(t.id)
	h.finished.append(t)
	h.log(fmt.Sprintf("worker #%d finished task #%d", workerId, t.id))
}

func (h *Herder) Kill(taskId int) error {
	h.active.m.Lock()
	defer h.active.m.Unlock()
	for i := range h.active.slice {
		if h.active.slice[i].id != taskId {
			continue
		}
		if h.active.slice[i].p != nil {
			return h.active.slice[i].p.kill()
		}
	}
	return errors.New("not found")
}

func (h *Herder) GetAllStates() []TaskState {
	q := h.GetQueue()
	a := h.GetActive()
	f := h.GetFinished()
	res := make([]TaskState, 0, len(q)+len(a)+len(f))
	res = append(res, q...)
	res = append(res, a...)
	res = append(res, f...)
	return res
}
func (h *Herder) GetQueue() []TaskState {
	return h.queue.getStates()
}
func (h *Herder) GetActive() []TaskState {
	return h.active.getStates()
}
func (h *Herder) GetFinished() []TaskState {
	return h.finished.getStates()
}
