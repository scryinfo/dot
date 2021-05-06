//
//Created by jayce on 9:13 2021/5/6
//
package data

import "testing"

func TestFindDots(t *testing.T) {

	var dirs []string
	dirs = append(dirs, `F:\go\src\github.com\scryinfo\dot\dots`)
	dirs = append(dirs, `F:\go\src\github.com\scryinfo\dot\dots\gindot`)
	data, ndir, err := FindDots(dirs)
	_ = data
	_ = ndir
	_ = err
}
