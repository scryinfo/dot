package main

import "github.com/scryinfo/dot/line"

func main() {
	l := line.New()
	l.ToLifer().Create(nil)

}

func add(l line.Line) {
	l.PreAdd()
}
