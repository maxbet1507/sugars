package sugars

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func ExamplePulse() {
	ctx, cancel := context.WithTimeout(context.Background(), 5500*time.Millisecond)
	defer cancel()

	pulse, stop := Pulse(ctx, 1*time.Second)
	defer stop()

	for pulse() {
		fmt.Println("hello")
	}
	fmt.Println("world")

	// manual stop for internal time.Ticker
	// stop()

	// Output:
	// hello
	// hello
	// hello
	// hello
	// hello
	// world
}

func TestPulseStopBeforeCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	pulse, stop := Pulse(ctx, 1*time.Second)
	defer stop()

	i := 0
	for pulse() {
		i++
	}

	stop()
	stop()
	cancel()

	if i != 1 {
		t.Fatal(i)
	}
}

func TestPulseStopAfterCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	pulse, stop := Pulse(ctx, 1*time.Second)
	defer stop()

	i := 0
	for pulse() {
		i++
	}

	cancel()
	stop()
	stop()

	if i != 1 {
		t.Fatal(i)
	}
}

func TestPulseCancelBeforePulse(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	pulse, stop := Pulse(ctx, 1*time.Second)
	defer stop()

	cancel()

	i := 0
	for pulse() {
		i++
	}

	stop()
	stop()

	if i != 0 {
		t.Fatal(i)
	}
}
