package herder

import (
	"log"
	"testing"
	"time"
)

func TestHerder(t *testing.T) {
	h := New(Config{
		MaxWorkersCount:     2,
		Logger:              log.Default(),
		DefaultMaxStdoutLen: 10,
		DefaultMaxStderrLen: 10,
	})
	var taskId int
	taskId = h.AddToQueue(TaskConfig{
		Command: "python",
		Args:    []string{"-u", "-c", "print(123,end='')"},
	})
	if taskId <= 0 {
		t.Error("taskId <= 0")
	}
	go h.Run()
	time.Sleep(2 * time.Second)
	states := h.GetAllStates()
	if len(states) != 1 {
		t.Error("len(states) != 1")
	}
	output := string(states[0].StdOut)
	expected := "123"
	if output != expected {
		t.Errorf("error on checking task stdout (wants \"%s\" found \"%s\")", expected, output)
	}
	err := h.Kill(taskId)
	if err != nil {
		t.Errorf("error in h.Kill(%d): %s", taskId, err.Error())
	}
}
