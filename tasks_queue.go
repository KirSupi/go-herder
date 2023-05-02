package herder

import (
	"errors"
	"sync"
)

type TasksQueue struct {
	len          int
	first        *tasksQueueItem
	last         *tasksQueueItem
	ch           chan *task
	nowInChannel *task
	m            sync.Mutex
}

func newTasksQueue() TasksQueue {
	return TasksQueue{
		ch: make(chan *task),
	}
}

type tasksQueueItem struct {
	task *task
	next *tasksQueueItem
}

func (q *TasksQueue) add(t *task) {
	q.m.Lock()
	item := &tasksQueueItem{
		task: t,
		next: nil,
	}
	if q.last != nil {
		q.last.next = item
		q.last = item
	} else {
		q.first = item
		q.last = item
	}
	q.len++
	if q.len == 1 {
		go q.run()
	}
	q.m.Unlock()
}

func (q *TasksQueue) remove(taskId int) error {
	q.m.Lock()
	defer q.m.Unlock()
	var prev, item *tasksQueueItem
	for prev, item = nil, q.first; item != nil; prev, item = item, item.next {
		if item.task.id == taskId {
			if prev == nil {
				q.first = item.next
				if q.len <= 2 {
					q.last = item.next
				}
			} else {
				prev.next = item.next
				if item.next == nil {
					q.last = prev
				}
			}
			q.len--
		}
	}
	return errors.New("process not found")
}

func (q *TasksQueue) clear() {
	q.m.Lock()
	defer q.m.Unlock()
	q.first = nil
	q.last = nil
	q.len = 0
}

func (q *TasksQueue) pop() *task {
	if q.len == 0 {
		return nil
	}
	q.m.Lock()
	defer q.m.Unlock()
	var item *tasksQueueItem
	if q.len == 1 {
		item = q.first
		q.first = nil
		q.last = nil
	} else {
		item = q.first
		q.first = q.first.next
	}
	q.len--
	return item.task
}

func (q *TasksQueue) run() {
	for t := q.pop(); t != nil; t = q.pop() {
		q.nowInChannel = t
		q.ch <- t
		q.nowInChannel = nil
	}
}

func (q *TasksQueue) iterate() chan *task {
	q.m.Lock()
	ch := make(chan *task)
	go func() {
		for i := q.first; i != nil; i = i.next {
			ch <- i.task
		}
		close(ch)
		q.m.Unlock()
	}()
	return ch
}

func (q *TasksQueue) getStates() []TaskState {
	q.m.Lock()
	defer q.m.Unlock()
	tasksCount := q.len
	if q.nowInChannel != nil {
		tasksCount++
	}
	res := make([]TaskState, 0, tasksCount)
	if q.nowInChannel != nil {
		res = append(res, q.nowInChannel.getState())
	}
	for i := q.first; i != nil; i = i.next {
		res = append(res, i.task.getState())
	}
	return res
}
