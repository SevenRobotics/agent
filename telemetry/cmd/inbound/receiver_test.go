package inbound

import (
	"context"
	"errors"
	"os"
	"strconv"
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
	if strings.Contains(source, "Published %s from queue=") {
		t.Fatal("receiver still contains verbose published-data log")
	}
}

func TestUniqueROSNodeNameKeepsConfiguredBaseAndAddsProcessSuffix(t *testing.T) {
	nodeName := uniqueROSNodeName("joystick_cmd_vel_publisher")

	if !strings.HasPrefix(nodeName, "joystick_cmd_vel_publisher_") {
		t.Fatalf("expected configured base prefix, got %q", nodeName)
	}
	if !strings.HasSuffix(nodeName, "_"+strconv.Itoa(os.Getpid())) {
		t.Fatalf("expected pid suffix, got %q", nodeName)
	}
}

func TestSanitizeROSNameTokenReplacesInvalidCharacters(t *testing.T) {
	got := sanitizeROSNameToken("robot-01.local / inbound")
	want := "robot_01_local_inbound"

	if got != want {
		t.Fatalf("sanitizeROSNameToken() = %q, want %q", got, want)
	}
}

func TestNewReceiverValidation(t *testing.T) {
	rmqConfig := config.RMQConfig{}
	rosConfig := config.RosNodeConfig{
		Address: "localhost:11311",
	}

	// Valid cases
	t.Run("valid twist", func(t *testing.T) {
		cfg := config.RMQInboundReceiverConfig{
			Name:        "test-twist",
			Queue:       "test-q",
			MessageType: "geometry_msgs/Twist",
			RosTopic:    "/cmd_vel",
		}
		_, err := NewReceiver(rmqConfig, rosConfig, cfg)
		if err != nil {
			t.Fatalf("unexpected error for twist message type: %v", err)
		}
	})

	t.Run("valid string", func(t *testing.T) {
		cfg := config.RMQInboundReceiverConfig{
			Name:        "test-string",
			Queue:       "test-q",
			MessageType: "std_msgs/String",
			RosTopic:    "/tasks",
		}
		r, err := NewReceiver(rmqConfig, rosConfig, cfg)
		if err != nil {
			t.Fatalf("unexpected error for string message type: %v", err)
		}
		if r.config.RosNodeName != "test-string_publisher" {
			t.Fatalf("expected RosNodeName %q, got %q", "test-string_publisher", r.config.RosNodeName)
		}
	})

	t.Run("wait_for_ros_subscriber defaults and bypasses", func(t *testing.T) {
		cfgDefault := config.RMQInboundReceiverConfig{
			Name:        "test-default-wait",
			Queue:       "test-q",
			MessageType: "std_msgs/String",
			RosTopic:    "/tasks",
		}
		rDefault, err := NewReceiver(rmqConfig, rosConfig, cfgDefault)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rDefault.config.WaitForRosSubscriber != "/tasks_node" {
			t.Fatalf("expected default WaitForRosSubscriber to be /tasks_node, got %q", rDefault.config.WaitForRosSubscriber)
		}

		cfgNone := config.RMQInboundReceiverConfig{
			Name:                 "test-none-wait",
			Queue:                "test-q",
			MessageType:          "std_msgs/String",
			RosTopic:             "/tasks",
			WaitForRosSubscriber: "none",
		}
		rNone, err := NewReceiver(rmqConfig, rosConfig, cfgNone)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rNone.config.WaitForRosSubscriber != "" {
			t.Fatalf("expected WaitForRosSubscriber to be empty for 'none', got %q", rNone.config.WaitForRosSubscriber)
		}

		cfgDisabled := config.RMQInboundReceiverConfig{
			Name:                 "test-disabled-wait",
			Queue:                "test-q",
			MessageType:          "std_msgs/String",
			RosTopic:             "/tasks",
			WaitForRosSubscriber: "disabled",
		}
		rDisabled, err := NewReceiver(rmqConfig, rosConfig, cfgDisabled)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rDisabled.config.WaitForRosSubscriber != "" {
			t.Fatalf("expected WaitForRosSubscriber to be empty for 'disabled', got %q", rDisabled.config.WaitForRosSubscriber)
		}
	})

	// Invalid case
	t.Run("invalid type", func(t *testing.T) {
		cfg := config.RMQInboundReceiverConfig{
			Name:        "test-invalid",
			Queue:       "test-q",
			MessageType: "invalid_msgs/Custom",
			RosTopic:    "/custom",
		}
		_, err := NewReceiver(rmqConfig, rosConfig, cfg)
		if err == nil {
			t.Fatal("expected error for unsupported message type, got nil")
		}
	})
}
