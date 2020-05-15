package field

import (
	"testing"
)

func (f *Field) TestCelSelectedTimesCount(t *testing.T) {
	result := f.CellSelectedTimesCount(hoge)
	expext := fuga
	if result != expext {
		// error‚ð‹L˜^
		// t.Error("\nresult: ", result, "\nexpext: ", expext)
	}

	t.Log("Test is finished!")
}

func (f *Field) TestActAgents(t *testing.T) {
	result := f.ActAgents(hoge)
	expext := fuga
	if result != expext {
		// error‚ð‹L˜^
	}
	t.Log("Test is finished")
}