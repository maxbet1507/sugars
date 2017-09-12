package ctxmain

type status struct {
	err  error
	ret  chan error
	once chan struct{}
}

func (s *status) Report(err error) {
	if _, ok := <-s.once; ok {
		s.ret <- err
		close(s.once)
		close(s.ret)
	}
}

func (s *status) Result() error {
	if err, ok := <-s.ret; ok {
		s.err = err
	}
	return s.err
}

// NewStatus -
func NewStatus() Status {
	ret := make(chan error, 1)
	once := make(chan struct{}, 1)
	once <- struct{}{}

	return &status{
		ret:  ret,
		once: once,
	}
}
