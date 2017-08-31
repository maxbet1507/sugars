package sugars

import (
	"context"
	"fmt"
	"testing"
)

func ExampleServiceFunc() {
	fn := func(ctx context.Context) {
		fmt.Println("running")
		defer fmt.Println("stopped")
		<-ctx.Done()
	}

	svc := ServiceFunc(context.Background(), fn)
	fmt.Println(1, svc.Running())

	svc.Start()
	fmt.Println(2, svc.Running())

	svc.Start()
	fmt.Println(3, svc.Running())

	svc.Stop()
	fmt.Println(4, svc.Running())

	svc.Stop()
	fmt.Println(5, svc.Running())

	svc.Start()
	fmt.Println(6, svc.Running())

	svc.Stop()
	fmt.Println(7, svc.Running())

	svc.Dispose()
	fmt.Println(8, svc.Running())

	// Output:
	// 1 false
	// running
	// 2 true
	// 3 true
	// stopped
	// 4 false
	// 5 false
	// running
	// 6 true
	// stopped
	// 7 false
	// 8 false
}

func TestServiceFuncContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	svc := ServiceFunc(ctx, func(ctx context.Context) {
		cancel()
		<-ctx.Done()
	})
	svc.Start()
	svc.Dispose()
}
