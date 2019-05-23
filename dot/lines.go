// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

var gDefaultline Line

//Return default Line，general program only has one line 
func GetDefaultLine() Line {
	return gDefaultline
}

//Set default line
func SetDefaultLine(line Line) {
	gDefaultline = line
}

