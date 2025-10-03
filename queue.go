package dorky

// Queue implements a generic queue
type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, 10),
	}
}

func (q *Queue[T]) enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) enqueueMultiple(items []T) {
	q.items = append(q.items, items...)
}

func (q *Queue[T]) dequeue() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}

	item := q.items[0]
	q.items[0] = zero // Clear reference to prevent memory leak
	q.items = q.items[1:]

	return item, true
}
