package netutils

type NetQueueItem struct {
	content []byte
}

type NetQueue struct {
	// the first item is the first in the queue aka. the one with the most priority
	queue []NetQueueItem
}

// NewQueue creates a new NetQueue instance.
func NewQueue() NetQueue {
	nq := NetQueue{}

	nq.queue = make([]NetQueueItem, 0)

	return nq
}

// NewItem creates a new NetQueueItem instance.
func NewItem(content []byte) NetQueueItem {
	nqi := NetQueueItem{}

	nqi.content = content

	return nqi
}

func (nqi NetQueueItem) Content() []byte {
	return nqi.content
}

// AddItem adds an item as last in the queue.
func (nq *NetQueue) AddItem(item NetQueueItem) {
	nq.queue = append(nq.queue, item)
}

// Next pops the first item and returns it.
func (nq *NetQueue) Next() NetQueueItem {
	// get the first item
	item := nq.queue[0]

	// remove the first item
	nq.queue = append(nq.queue[:0], nq.queue[1:]...)

	return item
}

// IsEmpty returns true if the queue is empty.
func (nq NetQueue) IsEmpty() bool {
	return len(nq.queue) == 0
}
