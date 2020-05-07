package utils

import "testing"

func TestUTILS(t *testing.T) {
	contains := Contains([]string{"a", "b", "c"}, "a")
	if contains == false {
		t.Error("utils fails contains")
	}
}
