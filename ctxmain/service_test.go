package ctxmain

import (
	"context"
	"fmt"
	"testing"
)

var (
	errMain = fmt.Errorf("test")
)

type WaitMain struct{}

func (s WaitMain) Main(ctx context.Context) error {
	<-ctx.Done()
	return errMain
}

func TestWaitMain(t *testing.T) {
	s := NewService(context.Background(), WaitMain{})

	if r := s.Running(); r {
		t.Fatal(r)
	}
	if r := s.Error(); r != nil {
		t.Fatal(r)
	}

	s.Start()

	if r := s.Running(); !r {
		t.Fatal(r)
	}
	if r := s.Error(); r != nil {
		t.Fatal(r)
	}

	s.Stop()

	if r := s.Running(); r {
		t.Fatal(r)
	}
	if r := s.Error(); r != errMain {
		t.Fatal(r)
	}

	// restart/restop

	s.Start()

	if r := s.Running(); !r {
		t.Fatal(r)
	}
	if r := s.Error(); r != nil {
		t.Fatal(r)
	}

	s.Start()

	if r := s.Running(); !r {
		t.Fatal(r)
	}
	if r := s.Error(); r != nil {
		t.Fatal(r)
	}

	s.Stop()

	if r := s.Running(); r {
		t.Fatal(r)
	}
	if r := s.Error(); r != errMain {
		t.Fatal(r)
	}

	s.Stop()

	if r := s.Running(); r {
		t.Fatal(r)
	}
	if r := s.Error(); r != errMain {
		t.Fatal(r)
	}

	s.Dispose()
}

type NowaitMain struct {
	cancel func()
}

func (s NowaitMain) Main(context.Context) error {
	defer s.cancel()
	return errMain
}

func TestNowaitMain(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	s := NewService(context.Background(), NowaitMain{cancel: cancel})

	s.Start()
	<-ctx.Done()

	if r := s.Running(); !r {
		t.Fatal(r)
	}
	if r := s.Error(); r != nil {
		t.Fatal(r)
	}

	s.Stop()

	if r := s.Running(); r {
		t.Fatal(r)
	}
	if r := s.Error(); r != errMain {
		t.Fatal(r)
	}

	s.Dispose()
}
