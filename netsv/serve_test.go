package netsv

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func Example() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	for i := 0; i < 3; i++ {
		go func() {
			conn, err := net.Dial("tcp", "127.0.0.1:10000")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
		}()
	}

	fn := func(ctx context.Context, conn net.Conn) {
		fmt.Println("serving")

		<-ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}

	if err := ServeFunc(ctx, "tcp", "127.0.0.1:10000", fn); err != nil {
		fmt.Println(err)
	}

	fmt.Println("graceful")
	// Output:
	// serving
	// serving
	// serving
	// graceful
}

func TestServeListenError(t *testing.T) {
	fn := func(ctx context.Context, conn net.Conn) {
	}

	if err := ServeFunc(context.Background(), "invalid", "address", fn); err == nil {
		t.Fatal(err)
	}
}
