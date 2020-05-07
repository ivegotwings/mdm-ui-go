package state

import "testing"

func TestSTATE(t *testing.T) {
	err := Connect(nil)
	if err != nil {
		t.Errorf("state failed to connect to %v", err)
	}
}
