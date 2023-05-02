package herder

import "testing"

func TestTasksQueue(t *testing.T) {
	q := newTasksQueue()
	q.add(&task{})
}
