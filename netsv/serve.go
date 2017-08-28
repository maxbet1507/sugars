package netsv

import (
	"context"
	"net"
	"strings"
)

type msgCloseAll struct {
}

type msgClose struct {
	Conn net.Conn
}

type msgAccept struct {
	Conn net.Conn
}

type manager struct {
	command  chan interface{}
	finished chan struct{}
}

func (s *manager) CloseAll() {
	s.command <- msgCloseAll{}
	<-s.finished
}

func (s *manager) Close(conn net.Conn) {
	s.command <- msgClose{Conn: conn}
}

func (s *manager) Accept(conn net.Conn) {
	s.command <- msgAccept{Conn: conn}
}

func netConnKey(conn net.Conn) string {
	return strings.Join([]string{
		conn.LocalAddr().String(),
		conn.RemoteAddr().String(),
	}, " ")
}

func (s *manager) Run() {
	s.command = make(chan interface{})
	defer close(s.command)

	s.finished = make(chan struct{})
	defer close(s.finished)

	conns := map[string]net.Conn{}
	finishing := false

	for msg := range s.command {
		switch msg.(type) {
		case msgAccept:
			conn := msg.(msgAccept).Conn
			key := netConnKey(conn)

			conns[key] = conn

		case msgClose:
			conn := msg.(msgClose).Conn
			key := netConnKey(conn)

			conn.Close()
			delete(conns, key)

		case msgCloseAll:
			for _, conn := range conns {
				conn.Close()
			}
			finishing = true
		}

		if finishing && len(conns) == 0 {
			return
		}
	}
}

// Handler -
type Handler interface {
	Serve(context.Context, net.Conn)
}

func filterAcceptError(err error) error {
	if operr, ok := err.(*net.OpError); ok {
		if operr.Err.Error() == "use of closed network connection" {
			err = nil
		}
	}
	return err
}

// Serve -
func Serve(ctx context.Context, network, addr string, handler Handler) error {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	manager := manager{}
	go manager.Run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			manager.CloseAll()
			return filterAcceptError(err)
		}
		manager.Accept(conn)

		go func() {
			handler.Serve(ctx, conn)
			manager.Close(conn)
		}()
	}
}

type prxhandler struct {
	fn func(context.Context, net.Conn)
}

func (s prxhandler) Serve(ctx context.Context, conn net.Conn) {
	s.fn(ctx, conn)
}

// ServeFunc -
func ServeFunc(ctx context.Context, network, addr string, fn func(context.Context, net.Conn)) error {
	return Serve(ctx, network, addr, prxhandler{fn: fn})
}
