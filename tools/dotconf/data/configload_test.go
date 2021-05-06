package data

import "testing"

func TestPathDeal(t *testing.T) {
	path, err := pathDeal("testdata")
	if err != nil {
		return
	}
	t.Logf("%#v", path)
}
