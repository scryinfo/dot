// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

var defaultLine Line

//Return default Line，general program only has one line
func GetDefaultLine() Line {
	return defaultLine
}

//Set default line
func SetDefaultLine(line Line) {
	defaultLine = line
}
