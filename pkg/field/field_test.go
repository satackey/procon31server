package field

import (
	"testing"
)

func TestCelSelectedTimesCount(t *testing.T) {
	result := CellSelectedTimesCount(hoge)
	expext := fugafuga
	if result != expext {
		t.Error("\nresult: ", result, "\nexpext: ", expext)
	}

	t.Log("Test is finished!")
}