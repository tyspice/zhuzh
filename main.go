package main

import "github.com/tyspice/zhuzh/internal/ui"

// set during build
var version string = "dev"

func main() {
	ui.Run()
}
