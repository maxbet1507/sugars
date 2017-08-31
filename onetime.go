package sugars

import (
	"sync"
)

// Onetime -
func Onetime(fn func()) func() {
	once := new(sync.Once)
	return func() {
		once.Do(fn)
	}
}

