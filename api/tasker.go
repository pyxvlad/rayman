package api

import (
	"gitlab.com/rayone121/rayman/pacman"
	"sync"
)

type Operation struct {
	Operation pacman.Operation `json:"operation,omitempty"`
	Packages  []pacman.Package `json:"packages,omitempty"`
	Err       error            `json:"error,omitempty"`
	ID        int              `json:"id,omitempty"`
}

type Tasker struct {
	currentLock      sync.RWMutex
	currentOperation Operation

	operationFeed chan Operation

	completedLock sync.RWMutex
	completed     []Operation

	counterLock sync.Mutex
	counter     int
}

func NewTasker() *Tasker {
	var t Tasker
	t.operationFeed = make(chan Operation, 10)
	t.counter = 1
	go t.executor()

	return &t
}

func (t *Tasker) executor() {
	for {
		select {
		case op := <-t.operationFeed:
			{
				t.currentLock.Lock()
				t.currentOperation = op
				t.currentLock.Unlock()
				op.Packages, op.Err = op.Operation.Execute()
				t.completedLock.Lock()
				t.completed = append(t.completed, op)
				t.completedLock.Unlock()
				t.currentLock.Lock()
				t.currentOperation = Operation{ID: 0}
				t.currentLock.Unlock()
			}
		default:
			break
		}
	}
}

func (t *Tasker) Schedule(op pacman.Operation) int {
	t.counterLock.Lock()
	id := t.counter
	t.counter++
	defer t.counterLock.Unlock()
	t.operationFeed <- Operation{ID: id, Operation: op}
	return id
}

func (t *Tasker) GetCompleted() []Operation {
	t.completedLock.RLock()
	tmp := make([]Operation, len(t.completed))
	copy(tmp, t.completed)
	t.completedLock.RUnlock()
	return tmp
}

func (t *Tasker) GetCurrent() Operation {
	t.currentLock.RLock()
	defer t.currentLock.RUnlock()
	return t.currentOperation
}
