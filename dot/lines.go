// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

var gDefaultline Line

//返回默认的Line， 一般的程序中只会有一个line
func GetDefaultLine() Line {
	return gDefaultline
}

//设置默认的line
func SetDefaultLine(line Line) {
	gDefaultline = line
}

