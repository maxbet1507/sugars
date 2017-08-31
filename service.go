package sugars

import (
	"context"
	"sync"
)

// ServiceHandler -
type ServiceHandler interface {
	Serve(context.Context)
}

// ServiceInstance -
type ServiceInstance interface {
	Stop()
	Start()
	Running() bool
	Dispose()
}

type msgServiceInstanceStart struct {
	done func()
}

type msgServiceInstanceStop struct {
	done func()
}

type serviceInstance struct {
	command  chan interface{}
	running  chan bool
	handler  ServiceHandler
	disposer func()
}

func (s *serviceInstance) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	s.command <- msgServiceInstanceStart{
		done: cancel,
	}

	<-ctx.Done()
}

func (s *serviceInstance) Stop() {
	ctx, cancel := context.WithCancel(context.Background())

	s.command <- msgServiceInstanceStop{
		done: cancel,
	}

	<-ctx.Done()
}

func (s *serviceInstance) Running() bool {
	return <-s.running
}

func (s *serviceInstance) main(ctx context.Context) {
	var cancel func()

	run := func() func() {
		subctx, cancel := context.WithCancel(ctx)
		runctx, running := context.WithCancel(context.Background())

		wg := new(sync.WaitGroup)
		wg.Add(1)

		go func() {
			defer wg.Done()
			running()
			s.handler.Serve(subctx)
		}()

		<-runctx.Done()

		return func() {
			cancel()
			wg.Wait()
		}
	}

	for {
		select {
		case <-ctx.Done():
			if cancel != nil {
				cancel()
			}
			return

		case msg := <-s.command:
			switch msg.(type) {
			case msgServiceInstanceStart:
				if cancel == nil {
					cancel = run()
				}
				msg.(msgServiceInstanceStart).done()
			case msgServiceInstanceStop:
				if cancel != nil {
					cancel()
					cancel = nil
				}
				msg.(msgServiceInstanceStop).done()
			}

		case s.running <- (cancel != nil):
		}
	}
}

func (s *serviceInstance) Dispose() {
	s.disposer()
}

// Service -
func Service(ctx context.Context, handler ServiceHandler) ServiceInstance {
	subctx, cancel := context.WithCancel(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	command := make(chan interface{})
	running := make(chan bool)

	r := &serviceInstance{
		command: command,
		running: running,
		handler: handler,
		disposer: Onetime(func() {
			defer close(command)
			defer close(running)
			cancel()
			wg.Wait()
		}),
	}

	go func() {
		defer wg.Done()
		r.main(subctx)
	}()

	return r
}

type proxyServiceHandler struct {
	fn func(context.Context)
}

func (s proxyServiceHandler) Serve(ctx context.Context) {
	s.fn(ctx)
}

// ServiceFunc -
func ServiceFunc(ctx context.Context, fn func(context.Context)) ServiceInstance {
	return Service(ctx, proxyServiceHandler{fn: fn})
}
