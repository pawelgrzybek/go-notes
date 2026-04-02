package notifier

import "sync"

type Notifier struct {
	mu          sync.RWMutex
	subscribers map[chan struct{}]struct{}
}

func New() *Notifier {
	return &Notifier{
		subscribers: make(map[chan struct{}]struct{}),
	}
}

func (n *Notifier) Subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	n.mu.Lock()
	n.subscribers[ch] = struct{}{}
	n.mu.Unlock()
	return ch
}

func (n *Notifier) Unsubscribe(ch chan struct{}) {
	n.mu.Lock()
	delete(n.subscribers, ch)
	n.mu.Unlock()
	close(ch)
}

func (n *Notifier) Notify() {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for ch := range n.subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
