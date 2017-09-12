package ctxmain

import (
	"fmt"
	"testing"
)

func TestStatus(t *testing.T) {
	err1 := fmt.Errorf("err1")
	err2 := fmt.Errorf("err1")

	s := NewStatus()

	s.Report(err1)
	if err := s.Result(); err != err1 {
		t.Fatal(err)
	}

	s.Report(err2)
	if err := s.Result(); err != err1 {
		t.Fatal(err)
	}
}
