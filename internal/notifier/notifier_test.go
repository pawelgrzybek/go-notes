package notifier_test

import (
	"testing"

	"github.com/pawelgrzybek/go-notes/internal/notifier"
)

func TestSubscribeUnsubscribe(t *testing.T) {
	n := notifier.New()
	ch := n.Subscribe()

	n.Notify()

	select {
	case <-ch:
	default:
		t.Fatal("expected notification on subscribed channel")
	}

	n.Unsubscribe(ch)

	_, ok := <-ch
	if ok {
		t.Fatal("expected channel to be closed after unsubscribe")
	}
}

func TestNotifyMultipleSubscribers(t *testing.T) {
	n := notifier.New()
	ch1 := n.Subscribe()
	ch2 := n.Subscribe()

	n.Notify()

	select {
	case <-ch1:
	default:
		t.Fatal("expected notification on ch1")
	}

	select {
	case <-ch2:
	default:
		t.Fatal("expected notification on ch2")
	}
}

func TestNotifyCoalescesWhenNotConsumed(t *testing.T) {
	n := notifier.New()
	ch := n.Subscribe()

	n.Notify()
	n.Notify()
	n.Notify()

	select {
	case <-ch:
	default:
		t.Fatal("expected at least one notification")
	}

	select {
	case <-ch:
		t.Fatal("expected only one buffered notification")
	default:
	}
}

func TestNotifyWithNoSubscribers(t *testing.T) {
	n := notifier.New()
	n.Notify() // should not panic
}
