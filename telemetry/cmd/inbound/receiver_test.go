package inbound

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"go_agent/config"
)

func newTestReceiver(check func(context.Context) (bool, error)) *Receiver {
	return &Receiver{
		config: config.RMQInboundReceiverConfig{
			Name:                 "test-receiver",
			RosTopic:             "/cmd_vel",
			WaitForRosSubscriber: "/cmd_vel_node",
		},
		rosSubscriberReady: check,
		rosPollInterval:    time.Millisecond,
	}
}

func TestWaitForROSSubscriberRetriesErrorsUntilReady(t *testing.T) {
	calls := 0
	receiver := newTestReceiver(func(context.Context) (bool, error) {
		calls++
		if calls < 3 {
			return false, errors.New("ros master unavailable")
		}
		return true, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := receiver.waitForROSSubscriber(ctx); err != nil {
		t.Fatalf("waitForROSSubscriber returned error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 readiness checks, got %d", calls)
	}
}

func TestWaitForROSSubscriberReturnsOnContextCancellation(t *testing.T) {
	receiver := newTestReceiver(func(context.Context) (bool, error) {
		return false, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := receiver.waitForROSSubscriber(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestWaitForROSSubscriberWaitsUntilSubscriberPresent(t *testing.T) {
	calls := 0
	receiver := newTestReceiver(func(context.Context) (bool, error) {
		calls++
		return calls >= 3, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := receiver.waitForROSSubscriber(ctx); err != nil {
		t.Fatalf("waitForROSSubscriber returned error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 readiness checks, got %d", calls)
	}
}

func TestMonitorROSSubscriberReportsLostSubscriber(t *testing.T) {
	receiver := newTestReceiver(func(context.Context) (bool, error) {
		return false, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case err := <-receiver.monitorROSSubscriber(ctx):
		if err == nil {
			t.Fatal("expected monitor to report an error")
		}
		if !strings.Contains(err.Error(), "not available") {
			t.Fatalf("expected missing subscriber error, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ROS monitor error")
	}
}

func TestReceiverDoesNotContainVerboseReceiveLog(t *testing.T) {
	src, err := os.ReadFile("receiver.go")
	if err != nil {
		t.Fatalf("read receiver source: %v", err)
	}

	source := string(src)
	if strings.Contains(source, "printTwist") {
		t.Fatal("receiver still contains printTwist helper or call")
	}
	if strings.Contains(source, "Received %s queue=") {
		t.Fatal("receiver still contains verbose received-data log")
	}
}
