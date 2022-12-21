package herder

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"text/tabwriter"
)

type HerderConfig struct {
}

type Herder struct {
	Processes []*Process
	c         *HerderConfig
	r         Repository
	m         *sync.Mutex
}

func NewHerder(c *HerderConfig, r Repository) *Herder {
	return &Herder{
		c: c,
		r: r,
		m: &sync.Mutex{},
	}
}

func (h *Herder) Init() error {
	h.m.Lock()
	defer h.m.Unlock()
	if h.c == nil {
		return errors.New("can't run Herder without *Config")
	}
	if h.r == nil {
		return errors.New("can't run Herder without Repository")
	}
	for pd := range h.r.IterProcesses() {
		h.Processes = append(h.Processes, &Process{
			ID:      pd.ID,
			Label:   pd.Label,
			Command: pd.Command,
			Params:  pd.Params,
		})
	}
	return nil
}

func (h *Herder) GetAllStates() (result []ProcessState, err error) {
	h.m.Lock()
	defer h.m.Unlock()
	result = make([]ProcessState, len(h.Processes))
	for i, p := range h.Processes {
		result[i] = p.getState()
	}
	return
}
func (h *Herder) RunAll() {
	h.m.Lock()
	defer h.m.Unlock()
	for _, p := range h.Processes {
		_ = p.run()
	}
}
func (h *Herder) RestartAll() {
	h.m.Lock()
	defer h.m.Unlock()
	var wg sync.WaitGroup
	for _, p := range h.Processes {
		wg.Add(1)
		go func(wg *sync.WaitGroup, p *Process) {
			_ = p.kill()
			_ = p.run()
			wg.Done()
		}(&wg, p)
	}
	wg.Wait()
}
func (h *Herder) KillAll() {
	h.m.Lock()
	defer h.m.Unlock()
	for _, p := range h.Processes {
		_ = p.kill()
	}
}

func (h *Herder) GetState(processID int) (ProcessState, error) {
	h.m.Lock()
	defer h.m.Unlock()
	for _, p := range h.Processes {
		if p.ID == processID {
			return p.getState(), nil
		}
	}
	return ProcessState{}, errorNoProcessID(processID)
}
func (h *Herder) Run(processID int) error {
	h.m.Lock()
	defer h.m.Unlock()
	for _, p := range h.Processes {
		if p.ID == processID {
			return p.run()
		}
	}
	return errorNoProcessID(processID)
}
func (h *Herder) Restart(processID int) error {
	h.m.Lock()
	defer h.m.Unlock()
	for _, p := range h.Processes {
		if p.ID == processID {
			_ = p.kill()
			_ = p.run()
			return nil
		}
	}
	return errorNoProcessID(processID)
}
func (h *Herder) Kill(processID int) error {
	h.m.Lock()
	defer h.m.Unlock()
	for _, p := range h.Processes {
		if p.ID == processID {
			return p.kill()
		}
	}
	return errorNoProcessID(processID)
}

func (h *Herder) CheckProcessExists(processID int) error {
	h.m.Lock()
	defer h.m.Unlock()
	for _, p := range h.Processes {
		if p.ID == processID {
			return nil
		}
	}
	return errorNoProcessID(processID)
}

func StringifyStates(states []ProcessState) string {
	needToPrintLabel := false
	for _, s := range states {
		if s.Label != nil {
			needToPrintLabel = true
			break
		}
	}

	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', tabwriter.Debug)
	if needToPrintLabel {
		_, _ = fmt.Fprint(w, "ID\tLabel\tActive\tOutput")
	} else {
		_, _ = fmt.Fprint(w, "ID\tActive\tOutput")
	}
	for _, s := range states {
		if needToPrintLabel {
			if s.Label != nil {
				_, _ = fmt.Fprintf(w, "\n%d\t%s\t%v\t%v", s.ID, *s.Label, s.Active, string(s.Output))
			} else {
				_, _ = fmt.Fprintf(w, "\n%d\t\t%v\t%v", s.ID, s.Active, string(s.Output))
			}
		} else {
			_, _ = fmt.Fprintf(w, "\n%d\t%v\t%v", s.ID, s.Active, string(s.Output))
		}
	}
	_ = w.Flush()
	return b.String()
}

func errorNoProcessID(id int) error {
	return errors.New("no process with id #" + strconv.Itoa(id))
}
