package ctxmain

import (
	"context"
	"sync"

	"github.com/maxbet1507/sugars"
)

type service struct {
	command chan interface{}
	cancel  func()
}

type messageStart struct {
	done func()
}

type messageStop struct {
	done func()
}

type messageError struct {
	done func(error)
}

type messageRunning struct {
	done func(bool)
}

func (s *service) Start() {
	ret := make(chan struct{})
	s.command <- messageStart{
		done: func() { close(ret) },
	}
	<-ret
}

func (s *service) Stop() {
	ret := make(chan struct{})
	s.command <- messageStop{
		done: func() { close(ret) },
	}
	<-ret
}

func (s *service) Error() error {
	ret := make(chan error)
	defer close(ret)
	s.command <- messageError{
		done: func(err error) { ret <- err },
	}
	return <-ret
}

func (s *service) Running() bool {
	ret := make(chan bool)
	defer close(ret)
	s.command <- messageRunning{
		done: func(f bool) { ret <- f },
	}
	return <-ret
}

func (s *service) Dispose() {
	s.cancel()
}

func (s *service) main(ctx context.Context, m Main) {
	var runerr error
	var running bool

	reterr := make(chan error, 1)
	defer close(reterr)

	stop := func() {}
	defer stop()

	start := func() {
		wg := new(sync.WaitGroup)
		wg.Add(1)

		subctx, cancel := context.WithCancel(ctx)

		runerr = nil
		running = true

		go func() {
			reterr <- m.Main(subctx)
			wg.Done()
		}()

		stop = sugars.Onetime(func() {
			cancel()
			wg.Wait()
			runerr = <-reterr
			running = false
		})
	}

	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-s.command:
			switch msg.(type) {
			case messageStart:
				msg := msg.(messageStart)
				stop()
				start()
				msg.done()

			case messageStop:
				msg := msg.(messageStop)
				stop()
				msg.done()

			case messageError:
				msg := msg.(messageError)
				msg.done(runerr)

			case messageRunning:
				msg := msg.(messageRunning)
				msg.done(running)
			}
		}
	}
}

// NewService -
func NewService(ctx context.Context, m Main) Service {
	subctx, cancel := context.WithCancel(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	r := &service{
		command: make(chan interface{}),
		cancel: sugars.Onetime(func() {
			cancel()
			wg.Wait()
		}),
	}
	go func() {
		defer close(r.command)
		defer wg.Done()
		r.main(subctx, m)
	}()

	return r
}
