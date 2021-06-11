package reconciliation

import "github.com/oslokommune/okctl/pkg/config/constant"

// Queue handles FIFO queue operations on Reconcilers
type Queue struct {
	reconcilers []Reconciler
	requeues    map[string]int
}

// Pop removes and returns the first element in the list
func (q *Queue) Pop() Reconciler {
	if len(q.reconcilers) == 0 {
		return nil
	}

	reconciler := q.reconcilers[0]
	q.reconcilers = q.reconcilers[1:]

	return reconciler
}

// Push adds a Reconciler to the back of the queue
func (q *Queue) Push(reconciler Reconciler) error {
	if q.requeues[reconciler.String()] == constant.DefaultMaxReconciliationRequeues {
		return ErrMaximumReconciliationRequeues
	}

	q.reconcilers = append(q.reconcilers, reconciler)

	q.requeues[reconciler.String()]++

	return nil
}

// NewQueue creates a new queue from a list of Reconcilers
func NewQueue(reconcilers []Reconciler) Queue {
	queue := Queue{}

	queue.reconcilers = make([]Reconciler, len(reconcilers))
	queue.requeues = make(map[string]int, len(reconcilers))

	copy(queue.reconcilers, reconcilers)

	return queue
}
