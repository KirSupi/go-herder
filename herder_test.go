package herder

import (
	"fmt"
	"testing"
)

func TestHerder(t *testing.T) {
	h := NewLimited(2)
	var taskId int
	taskId = h.AddToQueue("cmd_with_args", "arg1", "arg2")
	fmt.Println("added to queue", taskId)
	taskId = h.AddToQueue("cmd_without_args")
	fmt.Println("added to queue", taskId)
	fmt.Println(h.GetAllStates())
}
