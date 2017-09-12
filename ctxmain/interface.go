package ctxmain

import (
	"context"
)

// Main -
type Main interface {
	Main(context.Context) error
}

// Service -
type Service interface {
	Start()
	Stop()
	Error() error
	Running() bool
	Dispose()
}
