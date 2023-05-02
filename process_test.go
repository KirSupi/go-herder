package herder

import "testing"

func TestProcess(t *testing.T) {
	p := newProcess(newOutputBuffer(-1), newOutputBuffer(-1), "ver")
	p.run()
	err := p.kill()
	if err != nil {
		t.Errorf("p.kill() err: %s", err.Error())
	}
}
