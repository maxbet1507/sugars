package netsv

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:10000")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
	}()

	fn := func(ctx context.Context, conn net.Conn) {
		<-ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}

	if err := ServeFunc(ctx, "tcp", "127.0.0.1:10000", fn); err != nil {
		t.Fatal(err)
	}
}

func TestServeListenError(t *testing.T) {
	fn := func(ctx context.Context, conn net.Conn) {
	}

	if err := ServeFunc(context.Background(), "invalid", "address", fn); err == nil {
		t.Fatal(err)
	}
}
